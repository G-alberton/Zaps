let activeContact = null;
let toastTimer;
let modalAction = null;
let currentTarget = null;
let timerInterval;
let seconds = 0;
let pendingFile = null;

async function loadConversations(params) {
    try{
        const res = await fetch("http://localhost:8080/conversations");
        const data = await res.json();

        const list = document.querySelector('.contacts-list');
        list.innerHTML = "";

        data.forEach(conv => {
            const card = createContactCardHTML(
                conv.phone,
                conv.last_message,
                conv.conversation_id
            );

            list.appendChild(card);
        });
    } catch (err) {
        console.error("Error ao carregar conversas:", err);
    }
}

async function loadMessages(conversationID){
    try{
        const res = await fetch(`http://localhost:8080/messages?conversation_id=${conversationID}`);
        const messages = await res.json();

        const container = document.querySelector('.messages-container');
        container.innerHTML = "";

        messages.forEach(msg => {
            const row = document.createElement('div');

            const side = msg.direction === "outbound" ? "sent" : "received";

            row.classList.add(
                'message-row',
                side === 'sent' ? 'message-sent' : 'message-received'
            );

            row.innerHTML = createMessageHTML(
                msg.body,
                new Date(msg.timestamp).toLocaleTimeString().slice(0,5),
                msg.type,
                side
            );

            container.appendChild(row);
        });

        scrollToBottom();

        } catch (err) {
            console.error("Erro ao carregar mensagens:", err);
    }
}

function saveNewContact() {
    const nameInput = document.getElementById('new-contact-name');
    const numberInput = document.getElementById('new-contact-number');
    
    const name = nameInput.value.trim();
    const rawNumber = numberInput.value.trim();
    const id = cleaNumber(rawNumber);

    if(name && id) {
        const contactData = {
            id: id, 
            name: name, 
            lastMsg: "Nova conversa...",
            avatarUrl: "" 
        };

        const newCard = createContactCardHTML(contactData.name, contactData.lastMsg, contactData.id, contactData.avatarUrl);
        document.querySelector('.contacts-list').prepend(newCard);

        nameInput.value = '';
        numberInput.value = '';
        document.getElementById('add-contact-modal').style.display = 'none';
        showToast("Contato salvo!");

                
    } else {
        showToast("Por favor, insira um nome e um número válido.");
    }
}
function getInitials(name) {
    if (!name) return "??";
    return name
        .split(' ')
        .map(word => word[0])
        .slice(0, 2)
        .join('')
        .toUpperCase();
}

function convertToBase64(file){
    return new Promise((resolve, reject) => {
        const reader = new FileReader();
        reader.readAsDataURL(file);
        reader.onload = () => resolve(reader.result);
        reader.onerror = error => reject(error);
    });
}

function renderAndSave(content, time, type, side, caption) {
    const container = document.querySelector('.messages-container');
    const row = document.createElement('div');
    row.classList.add('message-row', `message-${side}`);
    
    row.innerHTML = createMessageHTML(content, time, type, side, caption);
    container.appendChild(row);
    
    if (activeContact) {
        const lastMsgText = type === 'text' ? content : (type === 'image' ? '📷 Foto' : '📂 Arquivo');
        updateLastMsgDisplay(activeContact, lastMsgText, type);
    }
    scrollToBottom();
}

function showToast(message){
    const toast = document.getElementById('toast');
    clearTimeout(toastTimer);
    toast.innerText = message;
    toast.classList.add('show');

    setTimeout(() => {
        toast.classList.remove('show');
    }, 2500);
}

function handleDelete(btn) {
    const row = btn.closest('.message-row');
    showCustomModal(
        "Apagar Mensagem?",
        "Esta ação não pode ser desfeita.",
        false,
        'DELETE_MSG',
        row,
        true
    );
}

function handleCopy(btn) {
    const bubble = btn.closest('.message-bubble');
    const textElement = bubble.querySelector('.message-text');
    if (textElement) {
        const textToCopy = textElement.innerText;
        navigator.clipboard.writeText(textToCopy).then(() => {
            showToast("Mensagem Copiada!");
            document.querySelectorAll('.message-menu').forEach(m => m.classList.remove('show'));
        });
    }
}

function showCustomModal(title, description, isInput, action, target, isDanger = false){
    const modal = document.getElementById('custom-modal');
    const input = document.getElementById('modal-input');
    const confirmBtn = document.getElementById('modal-confirm-btn');

    document.getElementById('modal-title').innerText = title;
    document.getElementById('modal-description').innerText = description;

    input.style.display = isInput ? 'block' : 'none';
    if(isInput) input.value = target.querySelector('.contact-name')?.innerText || "";

    confirmBtn.className = isDanger ? 'btn-danger' : '';
    confirmBtn.style.background = isDanger ? '#ff4d4d' : 'var(--primary-color)';

    modalAction = action;
    currentTarget = target;
    modal.classList.add('active');
}

function editContact(btn){
    showCustomModal("Editar Nome", "Digite o novo nome do contato:", true, 'EDIT', btn.closest('.contact-card'));
}

function deleteContact(btn){
    showCustomModal("Excluir Contato", "Tem certeza? Isso apagará o contato e as mensagens.", false, 'DELETE', btn.closest('.contact-card'), true);
}

function clearChat(btn) {
    showCustomModal("Limpar Conversa", "Deseja apagar todas as mensagens?", false, 'CLEAR', null, true);
}

function handleReply(btn) { console.log("Responder clicado"); }

document.addEventListener('DOMContentLoaded', () => {
    
    const elements = {
        chatInput: document.getElementById('chat-input'),
        sendBtn: document.querySelector('.send-button'),
        audioBtn: document.querySelector('.audio-button'),
        messagesContainer: document.querySelector('.messages-container'),
        searchInput: document.getElementById('message-search'),
        searchBtn: document.getElementById('search-button'),
        recordingStatus: document.getElementById('recording-status'),
        timerDisplay: document.getElementById('recording-timer'),
        fileInput: document.getElementById('file-input'),
        attachBtn: document.getElementById('attach-btn'),
        dropOverlay: document.getElementById('drop-zone-overlay'),
        chatHeaderName: document.querySelector('.chat-header-info strong'),
        contactSearch: document.getElementById('contact-search'),
        addContactBtn: document.getElementById('add-contact-btn'),
        addContactModal: document.getElementById('add-contact-modal'),
        saveContactBtn: document.getElementById('save-contact-btn'),
        closeContactModal: document.querySelector('.close-contact-modal'),
        contactsList: document.querySelector('.contacts-list'),
        toggleSidebarBtn: document.getElementById('toggle-sidebar'),
        sidebar: document.querySelector('.sidebar'),
        modal: document.getElementById('custom-modal'),
        confirmBtn: document.getElementById('modal-confirm-btn'),
        cancelBtn: document.getElementById('modal-cancel-btn'),
        modalInput: document.getElementById('modal-input'),
        numberInput: document.getElementById('new-contact-number'),
        cancelAddContactBtn: document.getElementById('cancel-add-contact-btn')
    };

    window.modalAction = null;
    window.currentTarget = null;
    window.handleDelete = handleDelete;
    window.handleCopy = handleCopy;
    window.editContact = editContact;
    window.deleteContact = deleteContact;
    window.clearChat = clearChat;
    window.handleReply = handleReply;

    if (elements.saveContactBtn) {
        elements.saveContactBtn.onclick = saveNewContact;
    }

    setupEventListeners(elements);
    loadConversations();
});

function setupEventListeners(elements) {
    const {
        chatInput, sendBtn, audioBtn, messagesContainer, searchInput, searchBtn,
        recordingStatus, timerDisplay, fileInput, attachBtn, dropOverlay,
        chatHeaderName, contactSearch, addContactBtn, addContactModal,
        saveContactBtn, closeContactModal, contactsList, toggleSidebarBtn,
        sidebar, modal, confirmBtn, cancelBtn, modalInput, numberInput, cancelAddContactBtn
    } = elements;

    numberInput.addEventListener('input', formatPhoneNumber);

    cancelBtn.onclick = () => modal.classList.remove('active');
    confirmBtn.onclick = handleModalConfirm;

    if (toggleSidebarBtn && sidebar) {
        toggleSidebarBtn.onclick = () => {
            sidebar.classList.toggle('hidden');
            toggleSidebarBtn.innerText = sidebar.classList.contains('hidden') ? "➡️" : "☰";
        };
    }

    cancelAddContactBtn.onclick = () => {
        document.getElementById('new-contact-name').value = '';
        document.getElementById('new-contact-number').value = '';
        addContactModal.style.display = 'none';
    };
    window.onclick = (event) => {
        if (event.target === addContactModal) {
            addContactModal.style.display = 'none';
        }
    };

    chatInput.addEventListener('input', handleChatInput);
    sendBtn.addEventListener('click', sendMessage);
    chatInput.addEventListener('keydown', (e) => {
        if (e.key === 'Enter' && !e.shiftKey) {
            e.preventDefault();
            sendMessage();
        }
    });

    audioBtn.addEventListener('mousedown', startRec);
    window.addEventListener('mouseup', stopRec);
    audioBtn.addEventListener('touchstart', startRec, { passive: false });
    audioBtn.addEventListener('touchend', stopRec);

    searchBtn.addEventListener('click', () => searchMessages(messagesContainer));

    attachBtn.addEventListener('click', () => fileInput.click());
    fileInput.addEventListener('change', (e) => handleFiles(e.target.files));

    window.addEventListener('dragover', (e) => {
        e.preventDefault();
        dropOverlay.classList.add('active');
    });
    window.addEventListener('drop', (e) => {
        e.preventDefault();
        dropOverlay.classList.remove('active');
        handleFiles(e.dataTransfer.files);
    });

    contactSearch.addEventListener('input', handleContactSearch);

    addContactBtn.onclick = () => addContactModal.style.display = 'flex';
    closeContactModal.onclick = () => addContactModal.style.display = 'none';
    const closeModalBtn = document.querySelector('.close-modal');
    if (closeModalBtn) {
        closeModalBtn.onclick = () => {
            const mediaModal = document.getElementById('media-preview-modal');
            mediaModal.style.display = 'none';
            mediaModal.classList.remove('active');
            pendingFile = null; 
        };
    }
    

    saveContactBtn.onclick = saveNewContact;

    document.addEventListener('click', handleMenuClicks);
    document.addEventListener('click', handleAudioPlay);
    
}

function formatPhoneNumber(e) {
    let value = e.target.value.replace(/\D/g, '');
    let formattedValue = '';

    if (value.length > 0) {
        formattedValue = '(' + value.substring(0, 2);
        if (value.length > 2) formattedValue += ') ' + value.substring(2, 7);
        if (value.length > 7) formattedValue += "-" + value.substring(7, 11);
    }
    e.target.value = formattedValue;
}

function handleModalConfirm() {
    const modal = document.getElementById('custom-modal');
    
    if (modalAction === 'EDIT') {
        const newName = document.getElementById('modal-input').value.trim();
        if (newName && currentTarget) {
            currentTarget.querySelector('.contact-name').innerText = newName;
            if (currentTarget.classList.contains('active')) {
                document.querySelector('.chat-header-info strong').innerText = newName;
            }
        }
    } else if (modalAction === 'DELETE') {
        handleContactDelete(currentTarget);
    } else if (modalAction === 'CLEAR') {
        document.querySelector('.messages-container').innerHTML = "";
    } else if (modalAction === 'DELETE_MSG'){
        if (currentTarget) currentTarget.remove();
    }

    modal.classList.remove('active');
    document.querySelectorAll('.contact-menu, .message-menu').forEach(m => m.classList.remove('show'));
}

function createContactCardHTML(name, lastMsg, id, avatarUrl) {
    const card = document.createElement('div');
    card.classList.add('contact-card');
    card.dataset.id = id;

    const initials = getInitials(name);
    
    const imgHTML = avatarUrl 
        ? `<img src="${avatarUrl}" onerror="this.style.display='none'">` 
        : '';

    card.innerHTML = `
        <div class="contact-avatar" data-initials="${initials}">
            ${imgHTML}
        </div>
        <div class="contact-info">
            <span class="contact-name">${name}</span>
            <span class="contact-last-msg">${lastMsg}</span>
        </div>
        <div class="contact-options-wrapper">
            <button class="contact-options-btn" type="button">⋮</button>
            <div class="contact-menu">
                <button onclick="editContact(this)">Editar Nome</button>
                <button onclick="clearChat(this)">Limpar Conversa</button>
                <hr>
                <button class="delete-contact-btn" onclick="deleteContact(this)">Excluir Contato</button>
            </div>
        </div>
    `;

    card.addEventListener('click', (e) => handleContactClick(card, name, id, e));
    
    const menuBtn = card.querySelector('.contact-options-btn');
    menuBtn.onclick = (e) => {
        e.stopPropagation();
        const menu = card.querySelector('.contact-menu');
        document.querySelectorAll('.contact-menu').forEach(m => {
            if (m !== menu) m.classList.remove('show');
        });
        menu.classList.toggle('show');
    };

    return card;
}

function handleContactClick(card, name, id, e) {
    if (e.target.closest('.contact-menu')) return;

    document.querySelectorAll('.contact-card').forEach(c => c.classList.remove('active'));
    card.classList.add('active');

    activeContact = id;
    
    document.querySelector('.chat-header-info strong').innerText = name;

    updateHeaderAvatar(card.dataset.avatar, name);
    
    const messagesContainer = document.querySelector('.messages-container');
    messagesContainer.innerHTML = "";

    if (window.innerWidth <= 768) {
        document.querySelector('.sidebar')?.classList.add('hidden');
    }

    activeContact = id;
    loadMessages(id);
}

function updateHeaderAvatar(url, name) {
    const headerAvatarContainer = document.querySelector('.chat-header .contact-avatar'); 
    if (!headerAvatarContainer) return;

    const initials = getInitials(name);
    headerAvatarContainer.setAttribute('data-initials', initials);
    
    if (url) {
        headerAvatarContainer.innerHTML = `<img src="${url}" onerror="this.style.display='none'">`;
    } else {
        headerAvatarContainer.innerHTML = '';
    }
}

function handleContactDelete(target) {
    const idToDelete = target.dataset.id;
    
    target.remove();
    document.querySelector('.messages-container').innerHTML = "";
    document.querySelector('.chat-header-info strong').innerText = "Selecione um contato";
    activeContact = null;
}

function cleaNumber(num) {
    return num.toString().replace(/\D/g, '');
}

function createMessageHTML(content = "", time = "", type = "text", side = "sent", caption = "") {
    const isSent = side === 'sent';
    const safeContent =  content ?? "";
    const safeCaption = caption ?? "";
    const checks = isSent ? '<span class="message-checks">✓✓</span>' : '';

    const playSymbol = '\u25B6\uFE0E';

    let mainContentHTML = "";

    if (type === 'image') {
        mainContentHTML = `
            <div class="message-media">
                <img src="${safeContent}" class="msg-img" style="width:100%; border-radius:8px; display:block;">
            </div>`;
    } else if (type === 'file'){
        mainContentHTML = `
            <div class="file-wrapper" style="background:rgba(0,0,0,0.1); padding:10px; border-radius:8px; display:flex; align-items:center; gap:10px; color:inherit;">
                <span>📂</span> <small style="word-break:break-all;">${safeContent}</small>
            </div>
        `
    } else if (type === 'audio') {
        mainContentHTML = `
            <div class="audio-player-container">
                <button class="audio-play-btn">${playSymbol}</button>
                <div class="audio-controls">
                    <div class="audio-waveform"><div class="audio-progress"></div></div>
                    <div class="audio-meta"><span class="audio-duration">${safeContent}</span></div>
                </div>
            </div>`;
    } else {
        mainContentHTML = `<div class="message-text">${safeContent}</div>`;
    }

    const captionHTML = (safeCaption && type !== 'text') 
        ? `<div class="message-caption" style="margin-top:8px; font-size:0.95rem;">${safeCaption}</div>` 
        : "";

    return `
        <div class="message-bubble ${type === 'audio' ? 'audio-bubble' : ''}">
            <div class="message-options-btn"><span>&#9013;</span></div>
            <div class="message-menu">
                ${type !== 'audio' ? '<button onclick="handleReply(this)">Responder</button>' : ''}
                ${type !== 'audio' ? '<button onclick="handleCopy(this)">Copiar</button>' : ''}
                <button class="delete-btn" onclick="handleDelete(this)">Apagar</button>
            </div>
            
            ${mainContentHTML}
            ${captionHTML}
            
            <span class="message-time">${time} ${checks}</span>
        </div>
    `;
    
}

async function sendMessage() {
    const chatInput = document.getElementById('chat-input');
    const text = chatInput.value.trim();

    if (!activeContact) {
        alert("Selecione um contato");
        return;
    }

    const now = new Date();
    const time = now.getHours().toString().padStart(2, '0') + ":" + now.getMinutes().toString().padStart(2, '0')

    if (pendingFile) {
        const type = pendingFile.type.startsWith('image/') ? 'image' : 'file';
        const content = msg.type === "image" ? msg.media_id : msg.boody;
        
        renderAndSave(content, time, type, 'sent', text);

        const formData = new FormData();
        formData.append("file", pendingFile);
        formData.append("to", activeContact);
        formData.append("caption",text);

        await fetch("http://localhost:8080/send-media", {
            method: "POST",
            body: formData
        })

        pendingFile = null;
        chatInput.placeholder = "Digite uma Mensagem";
        return;
    }

    if (text !== "") {
        try{
            console.log("Enviando mensagem para:", activeContact);

            await fetch("http://localhost:8080/send-message", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json"
                },
                body: JSON.stringify({
                    to: activeContact,
                    body: text
                })
            });

            console.log("Mensagem enviada com sucesso");

            await loadMessages(activeContact)

            chatInput.value = "";
            chatInput.style.height = 'auto';
            toggleInputButtons("");

        } catch (err) {
            console.error("Erro ao enviar mensagem:", err);
        }
    }
}

/*function sendMessage() {
    const chatInput = document.getElementById('chat-input');
    const text = chatInput.value.trim();
    const now = new Date();
    const time = now.getHours().toString().padStart(2, '0') + ":" + now.getMinutes().toString().padStart(2, '0');

    if (pendingFile) {
        const type = pendingFile.type.startsWith('image/') ? 'image' : 'file';
        const content = type === 'image' ? pendingFile.base64 : pendingFile.name;

        renderAndSave(content, time, type, 'sent', text);

        pendingFile = null;
        chatInput.placeholder = "Digite uma Mensagem";
    } 

    else if (text !== "") {
        renderAndSave(text, time, 'text', 'sent');
        
        chatInput.value = "";
        chatInput.style.height = 'auto';
        toggleInputButtons("");

        setTimeout(() => receiveMessage("Resposta automática", "text"), 1500);
        setTimeout(() => receiveMessage("0:05", "audio"), 3000);
    }
}*/

function receiveMessage(content, type = 'text') {
    const now = new Date();
    const time = now.getHours() + ":" + now.getMinutes().toString().padStart(2, '0');
    const messagesContainer = document.querySelector('.messages-container');
    
    const messageRow = document.createElement('div');
    messageRow.classList.add('message-row', 'message-received');
    messageRow.innerHTML = createMessageHTML(content, time, type);
    
    messagesContainer.appendChild(messageRow);

    if (activeContact) {
        updateLastMsgDisplay(activeContact, content, type);
    }

    scrollToBottom();
}

function startTimer() {
    seconds = 0;
    document.getElementById('recording-timer').innerText = "00:00";
    timerInterval = setInterval(() => {
        seconds++;
        const mins = Math.floor(seconds / 60).toString().padStart(2, '0');
        const secs = (seconds % 60).toString().padStart(2, '0');
        document.getElementById('recording-timer').innerText = `${mins}:${secs}`;
    }, 1000);
}

function startRec(e){
    if (e.cancelable) e.preventDefault();
    const audioBtn = document.querySelector('.audio-button');
    audioBtn.classList.add('recording');
    audioBtn.innerHTML = "🛑";
    document.getElementById('chat-input').style.visibility = 'hidden';
    document.getElementById('recording-status').setAttribute('style', 'display: flex !important');
    startTimer();
}

function stopRec() {
    const audioBtn = document.querySelector('.audio-button');
    if (!audioBtn.classList.contains('recording')) return;

    audioBtn.classList.remove('recording');
    audioBtn.innerHTML = "🎙️";
    document.getElementById('chat-input').style.visibility = 'visible';
    document.getElementById('recording-status').setAttribute('style', 'display: none !important');

    const finalTime = document.getElementById('recording-timer').innerText;
    clearInterval(timerInterval);

    if (seconds > 0) {
        const now = new Date();
        const time = now.getHours() + ":" + now.getMinutes().toString().padStart(2, '0');
        const messagesContainer = document.querySelector('.messages-container');
        
        const messageRow = document.createElement('div');
        messageRow.classList.add('message-row', 'message-sent'); 
        messageRow.innerHTML = createMessageHTML(finalTime, time, 'audio');
        
        messagesContainer.appendChild(messageRow);
        if (activeContact){
            updateLastMsgDisplay(activeContact, finalTime, 'audio');
        }
        scrollToBottom();
    }
}

function toggleInputButtons(text) {
    const hasText = text.length > 0;
    const audioBtn = document.querySelector('.audio-button');
    const sendBtn = document.querySelector('.send-button');
    audioBtn.style.display = hasText ? 'none' : 'flex';
    sendBtn.style.display = hasText ? 'flex' : 'none';
}

function scrollToBottom() {
    const messagesContainer = document.querySelector('.messages-container');
    messagesContainer.scrollTop = messagesContainer.scrollHeight;
}

function handleChatInput(e) {
    e.target.style.height = 'auto'; 
    e.target.style.height = (e.target.scrollHeight) + 'px'; 
    toggleInputButtons(e.target.value.trim());
}

function searchMessages(messagesContainer) {
    const term = document.getElementById('message-search').value.toLowerCase().trim();
    const messages = messagesContainer.querySelectorAll('.message-row');
    
    messages.forEach(row => {
        const text = row.querySelector('.message-text')?.innerText.toLowerCase() || "";
        const isAudio = row.querySelector('.audio-duration')?.innerText.toLowerCase() || "";
        
        if (text.includes(term) || isAudio.includes(term) || term === "") {
            row.style.display = 'flex';
        } else {
            row.style.display = 'none';
        }
    });
}

function handleContactSearch(e) {
    const term = e.target.value.toLowerCase();
    document.querySelectorAll('.contact-card').forEach(card => {
        const name = card.querySelector('.contact-name').innerText.toLowerCase();
        if (name.includes(term)) {
            card.style.display = "flex";
        } else {
            card.style.display = 'none';
        }
    });
}

async function handleFiles(files) {
    const file = files[0]; 
    if (!file || !activeContact) return;

    const previewModal = document.getElementById('media-preview-modal');
    const previewContainer = document.getElementById('preview-container');
    const captionInput = document.getElementById('media-caption');

    previewContainer.innerHTML = "";
    captionInput.value = "";
    pendingFile = file;

    if (file.type.startsWith('image/')) {
        const base64 = await convertToBase64(file);
        previewContainer.innerHTML = `<img src="${base64}" style="max-width:100%; max-height:300px; border-radius:8px; display:block; margin: 0 auto;">`;
        pendingFile.base64 = base64; 
    } else {
        previewContainer.innerHTML = `
            <div style="text-align:center; padding:20px;">
                <span style="font-size:50px;">📄</span>
                <p style="margin-top:10px; word-break:break-all;">${file.name}</p>
            </div>`;
    }
    previewModal.style.display = 'flex'; 
    previewModal.classList.add('active');
}

document.getElementById('confirm-send-btn').onclick = () => {
    const caption = document.getElementById('media-caption').value.trim();
    const now = new Date(); 
    const time = now.getHours().toString().padStart(2, '0') + ":" + now.getMinutes().toString().padStart(2, '0');

    const type = pendingFile.type.startsWith('image/') ? 'image' : 'file';
    const content = type === 'image' ? pendingFile.base64 : pendingFile.name;
    
    renderAndSave(content, time, type, 'sent', caption);

    document.getElementById('media-preview-modal').style.display = 'none';
    document.getElementById('media-preview-modal').classList.remove('active');
    pendingFile = null; 
};

function sendMediaMessage(content, time, type) {
    const messagesContainer = document.querySelector('.messages-container');
    const row = document.createElement('div');
    row.classList.add('message-row', 'message-sent');
    row.innerHTML = createMessageHTML(content, time, type);
    messagesContainer.appendChild(row);
    scrollToBottom();
}

function handleMenuClicks(e) {
    const optionsBtn = e.target.closest('.message-options-btn');
    const menu = optionsBtn ? optionsBtn.nextElementSibling : null;

    if (optionsBtn) {
        document.querySelectorAll('.message-menu').forEach(m => {
            if (m !== menu) m.classList.remove('show');
        });
        menu.classList.toggle('show');
    } 
    else if (!e.target.closest('.message-menu')) {
        document.querySelectorAll('.message-menu').forEach(m => m.classList.remove('show'));
    }
}

function handleAudioPlay(e) {
    const playBtn = e.target.closest('.audio-play-btn');
    if (!playBtn) return;

    const container = playBtn.closest('.audio-player-container');
    const progressBar = container.querySelector('.audio-progress');
    const durationText = container.querySelector('.audio-duration').innerHTML;

    const parts = durationText.split(':');
    const totalSeconds = parseInt(parts[0]) * 60 + parseInt(parts[1]);

    if (playBtn.innerText === "▶") {
        playBtn.innerText = "⏸";

        let currentPercent = 0;
        const intervalTime = 100;
        const increment = (100 / (totalSeconds * 1000)) * intervalTime;

        if (playBtn.dataset.intervalId) clearInterval(playBtn.dataset.intervalId);

        const animation = setInterval(() => {
            currentPercent += increment;
            if (currentPercent >= 100) {
                currentPercent = 100;
                clearInterval(animation);
                playBtn.innerText = "▶";
                setTimeout(() => {progressBar.style.width = "0%";}, 500);
            }
            progressBar.style.width = currentPercent + "%";
        }, intervalTime);

        playBtn.dataset.intervalId = animation;
    } else {
        playBtn.innerText = "▶";
        if (playBtn.dataset.intervalId) clearInterval(parseInt(playBtn.dataset.intervalId));
    }
}

function updateLastMsgDisplay(contactId, text, type) {
    const card = document.querySelector(`.contact-card[data-id="${contactId}"]`);
    const contactsList = document.querySelector('.contacts-list');
    
    if (card) {
        const lastMsgSpan = card.querySelector('.contact-last-msg');
        const displaySafeText = type === 'audio' ? '🎙️ Áudio' : text;
        lastMsgSpan.innerText = displaySafeText;
        contactsList.prepend(card);
    }
}