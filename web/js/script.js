let chatHistories = JSON.parse(localStorage.getItem('chatHistories')) || {};
let savedContacts = JSON.parse(localStorage.getItem('savedContacts')) || [];
let activeContact = null;
let toastTimer;
let modalAction = null;
let currentTarget = null;
let timerInterval;
let seconds = 0;

function saveToLocalStorage(contactName, content, time, type, side){
    if (!chatHistories[contactName]){
        chatHistories[contactName] = [];
    }

    chatHistories[contactName].push({
        content: content,
        time: time,
        type: type,
        side: side
    });

    localStorage.setItem('chatHistories', JSON.stringify(chatHistories));
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

    if (savedContacts.length > 0){
        savedContacts.forEach(contact => {
            const newCard = createContactCardHTML(contact.name, contact.lastMsg || "Nova conversa...", contact.id);
            elements.contactsList.appendChild(newCard);
        });
    }

    setupEventListeners(elements);

    window.modalAction = null;
    window.currentTarget = null;
    window.handleDelete = handleDelete;
    window.handleCopy = handleCopy;
    window.editContact = editContact;
    window.deleteContact = deleteContact;
    window.clearChat = clearChat;
    window.handleReply = handleReply;
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

function createContactCardHTML(name, lastMsg, id) {
    const card = document.createElement('div');
    card.classList.add('contact-card');
    card.dataset.id = id;

    card.innerHTML = `
        <div class="contact-avatar"></div>
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
    const messagesContainer = document.querySelector('.messages-container');
    messagesContainer.innerHTML = "";

    if (chatHistories[id]) {
        chatHistories[id].forEach(msg => {
            const row = document.createElement('div');
            row.classList.add('message-row', msg.side === 'sent' ? 'message-sent' : 'message-received');
            row.innerHTML = createMessageHTML(msg.content, msg.time, msg.type);
            messagesContainer.appendChild(row);
        });
        scrollToBottom();
    }

    if (window.innerWidth <= 768) {
        document.querySelector('.sidebar')?.classList.add('hidden');
    }
}

function handleContactDelete(target) {
    const idToDelete = target.dataset.id;
    
    savedContacts = savedContacts.filter(c => c.id !== idToDelete);
    delete chatHistories[idToDelete];
    
    localStorage.setItem('savedContacts', JSON.stringify(savedContacts));
    localStorage.setItem('chatHistories', JSON.stringify(chatHistories));
    
    target.remove();
    document.querySelector('.messages-container').innerHTML = "";
    document.querySelector('.chat-header-info strong').innerText = "Selecione um contato";
    activeContact = null;
}

function cleaNumber(num) {
    return num.toString().replace(/\D/g, '');
}

function saveNewContact() {
    const nameInput = document.getElementById('new-contact-name');
    const numberInput = document.getElementById('new-contact-number');
    const name = nameInput.value.trim();
    const rawNumber = numberInput.value.trim();
    const id = cleaNumber(rawNumber);

    if(name && id) {
        const exist = savedContacts.some(c => c.id === id);
        if (exist) {
            showToast("Este número já está cadastrado!");
            return;
        }

        const newCard = createContactCardHTML(name, "Nova conversa...", id);
        document.querySelector('.contacts-list').prepend(newCard);

        savedContacts.unshift({id: id, name: name, lastMsg: "Nova conversa..."});
        localStorage.setItem('savedContacts', JSON.stringify(savedContacts));

        nameInput.value = '';
        numberInput.value = '';
        document.getElementById('add-contact-modal').style.display = 'none';
        showToast("Contato salvo!");
    } else {
        showToast("Por favor, insira um nome e um número válido.");
    }
}

function createMessageHTML(content, time, type = 'text') {
    const isAudio = type === 'audio';
    const playSymbol = '\u25B6\uFE0E'; 
    const pauseSymbol = '\u23F8\uFE0E'; 

    return `
        <div class="message-bubble ${isAudio ? 'audio-bubble' : ''}">
            <div class="message-options-btn"><span>&#9013;</span></div>
            <div class="message-menu">
                ${!isAudio ? '<button onclick="handleReply(this)">Responder</button>' : ''}
                ${!isAudio ? '<button onclick="handleCopy(this)">Copiar</button>' : ''}
                <button class="delete-btn" onclick="handleDelete(this)">Apagar</button>
            </div>
            ${isAudio ? `
                <div class="audio-player-container">
                    <button class="audio-play-btn">${playSymbol}</button>
                    <div class="audio-controls">
                        <div class="audio-waveform"><div class="audio-progress"></div></div>
                        <div class="audio-meta"><span class="audio-duration">${content}</span></div>
                    </div>
                </div>
            ` : `<div class="message-text">${content}</div>`}
            <span class="message-time">${time}</span>
        </div>
    `;
}

function sendMessage() {
    const chatInput = document.getElementById('chat-input');
    const text = chatInput.value.trim();
    if (text !== "") {
        const now = new Date();
        const time = now.getHours() + ":" + now.getMinutes().toString().padStart(2, '0');
        const messagesContainer = document.querySelector('.messages-container');
        
        const messageRow = document.createElement('div');
        messageRow.classList.add('message-row', 'message-sent');
        messageRow.innerHTML = createMessageHTML(text, time, 'text');
        
        messagesContainer.appendChild(messageRow);
        if (activeContact) {
            saveToLocalStorage(activeContact, text, time, 'text', 'sent');
            updateLastMsgDisplay(activeContact, text, 'text');
        }
        chatInput.value = "";
        chatInput.style.height = 'auto';
        toggleInputButtons("");
        scrollToBottom();

        setTimeout(() => receiveMessage("Resposta automática", "text"), 1500);
        setTimeout(() => receiveMessage("0:05", "audio"), 3000);
    }
}

function receiveMessage(content, type = 'text') {
    const now = new Date();
    const time = now.getHours() + ":" + now.getMinutes().toString().padStart(2, '0');
    const messagesContainer = document.querySelector('.messages-container');
    
    const messageRow = document.createElement('div');
    messageRow.classList.add('message-row', 'message-received');
    messageRow.innerHTML = createMessageHTML(content, time, type);
    
    messagesContainer.appendChild(messageRow);

    if (activeContact) {
        saveToLocalStorage(activeContact, content, time, type, 'received');
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
            saveToLocalStorage(activeContact, finalTime, time, 'audio', 'sent');
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

function handleFiles(files) {
    Array.from(files).forEach(file => {
        const isImage = file.type.startsWith('image/');
        const time = new Date().getHours() + ":" + new Date().getMinutes().toString().padStart(2, '0');

        if (isImage){
            const imageURL = URL.createObjectURL(file);
            sendMediaMessage(imageURL, time, 'image');
        } else { 
            sendMediaMessage(file.name, time, 'file');
        }
    });
}

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

        const contactIndex = savedContacts.findIndex(c => c.id === contactId);
        if (contactIndex !== -1){
            savedContacts[contactIndex].lastMsg = displaySafeText;
            localStorage.setItem('savedContacts', JSON.stringify(savedContacts));
        }
    }
}