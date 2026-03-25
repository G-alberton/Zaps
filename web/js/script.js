function handleDelete(btn) {
    if (confirm("Deseja apagar esta mensagem?")) {
        btn.closest('.message-row').remove();
    }
}

function handleCopy(btn) {
    const bubble = btn.closest('.message-bubble');
    const textElement = bubble.querySelector('.message-text');
    if (textElement) {
        const textToCopy = textElement.innerText;
        navigator.clipboard.writeText(textToCopy).then(() => {
            alert("Copiado!");
            bubble.querySelector('.message-menu').classList.remove('show');
        });
    }
}

let modalAction = null;
let currentTarget = null;

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


/*Aqui é o DOM*/ 
document.addEventListener('DOMContentLoaded', () => {
    const chatInput = document.getElementById('chat-input');
    const sendBtn = document.querySelector('.send-button');
    const audioBtn = document.querySelector('.audio-button');
    const messagesContainer = document.querySelector('.messages-container');
    const searchInput = document.getElementById('message-search');
    const searchBtn = document.getElementById('search-button');
    const recordingStatus = document.getElementById('recording-status');
    const timerDisplay = document.getElementById('recording-timer');
    const fileInput = document.getElementById('file-input');
    const attachBnt = document.getElementById('attach-btn');
    const dropOverlay = document.getElementById('drop-zone-overlay');
    const contactCard = document.querySelectorAll('.contact-card');
    const chatHeaderName = document.querySelector('.chat-header-info strong');
    const contactSearch = document.getElementById('contact-search');
    const addContactBtn = document.getElementById('add-contact-btn');
    const addContactModal = document.getElementById('add-contact-modal');
    const saveContactBtn = document.getElementById('save-contact-btn');
    const closeContactModal = document.querySelector('.close-contact-modal');
    const contactsList = document.querySelector('.contacts-list');
    const toggleSidebarBtn = document.getElementById('toggle-sidebar');
    const sidebar = document.querySelector('.sidebar');
    const modal = document.getElementById('custom-modal');
    const confirmBtn = document.getElementById('modal-confirm-btn');
    const cancelBtn = document.getElementById('modal-cancel-btn');
    const input = document.getElementById('modal-input');

    cancelBtn.onclick = () => modal.classList.remove('active');

    confirmBtn.onclick = () => {
        if (modalAction === 'EDIT') {
            const newName = input.value.trim();
            if (newName) {
                currentTarget.querySelector('.contact-name').innerText = newName;
                if (currentTarget.classList.contains('active')) {
                    document.querySelector('.chat-header-info strong').innerText = newName;
                }
            }
        } 
        else if (modalAction === 'DELETE') {
            currentTarget.remove();
            document.querySelector('.messages-container').innerHTML = "";
            document.querySelector('.chat-header-info strong').innerText = "Selecione um contato";
        } 
        else if (modalAction === 'CLEAR') {
            document.querySelector('.messages-container').innerHTML = "";
        }

        modal.classList.remove('active');
        document.querySelectorAll('.contact-menu').forEach(m => m.classList.remove('show'));
    };
    
    let timerInterval;
    let seconds = 0;

    if (toggleSidebarBtn && sidebar) {
    toggleSidebarBtn.onclick = () => {
        sidebar.classList.toggle('hidden');
        
        if (sidebar.classList.contains('hidden')) {
            toggleSidebarBtn.innerText = "➡️"; 
        } else {
            toggleSidebarBtn.innerText = "☰"; 
        }
    };
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
        const text = chatInput.value.trim();
        if (text !== "") {
            const now = new Date();
            const time = now.getHours() + ":" + now.getMinutes().toString().padStart(2, '0');
            
            const messageRow = document.createElement('div');
            messageRow.classList.add('message-row', 'message-sent');
            messageRow.innerHTML = createMessageHTML(text, time, 'text');
            
            messagesContainer.appendChild(messageRow);
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
        
        const messageRow = document.createElement('div');
        messageRow.classList.add('message-row', 'message-received');
        messageRow.innerHTML = createMessageHTML(content, time, type);
        
        messagesContainer.appendChild(messageRow);
        scrollToBottom();
    }

    function startTimer() {
        seconds = 0;
        timerDisplay.innerText = "00:00";
        timerInterval = setInterval(() => {
            seconds++;
            const mins = Math.floor(seconds / 60).toString().padStart(2, '0');
            const secs = (seconds % 60).toString().padStart(2, '0');
            timerDisplay.innerText = `${mins}:${secs}`;
        }, 1000);
    }

    function startRec(e){
        if (e.cancelable) e.preventDefault();
        audioBtn.classList.add('recording');
        audioBtn.innerHTML = "🛑";
        chatInput.style.visibility = 'hidden';
        recordingStatus.setAttribute('style', 'display: flex !important');
        startTimer();
    }

    function stopRec() {
    if (!audioBtn.classList.contains('recording')) return;

    audioBtn.classList.remove('recording');
    audioBtn.innerHTML = "🎙️";
    chatInput.style.visibility = 'visible';
    recordingStatus.setAttribute('style', 'display: none !important');

    const finalTime = timerDisplay.innerText;
    clearInterval(timerInterval);

    
    if (seconds > 0) {
        const now = new Date();
        const time = now.getHours() + ":" + now.getMinutes().toString().padStart(2, '0');
        
        const messageRow = document.createElement('div');
        messageRow.classList.add('message-row', 'message-sent'); 

        
        messageRow.innerHTML = createMessageHTML(finalTime, time, 'audio');
        
        messagesContainer.appendChild(messageRow);
        scrollToBottom();
    }
}

    audioBtn.addEventListener('mousedown', startRec);
    window.addEventListener('mouseup', stopRec);

    audioBtn.addEventListener('touchstart', startRec, {passive: false});
    audioBtn.addEventListener('touchend', stopRec);

    function toggleInputButtons(text) {
        const hasText = text.length > 0;
        audioBtn.style.display = hasText ? 'none' : 'flex';
        sendBtn.style.display = hasText ? 'flex' : 'none';
    }

    function scrollToBottom() {
        messagesContainer.scrollTop = messagesContainer.scrollHeight;
    }

    chatInput.addEventListener('input', function() {
    this.style.height = 'auto'; 
    this.style.height = (this.scrollHeight) + 'px'; 
    toggleInputButtons(this.value.trim());
});

    sendBtn.addEventListener('click', sendMessage);
    chatInput.addEventListener('keydown', (e) => {
        if (e.key === 'Enter' && !e.shiftKey) {
            e.preventDefault();
            sendMessage();
        }
    });

    searchBtn.addEventListener('click', () => {
    const term = searchInput.value.toLowerCase().trim();
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
});

    document.addEventListener('click', (e) => {
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
});

    document.addEventListener('click', (e) => {
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
            clearInterval(parseInt(playBtn.dataset.intervalId));
        }
    });

    attachBnt.addEventListener('click', () => fileInput.click());

    fileInput.addEventListener('change', (e) => {
        handleFiles(e.target.files);
    });

    window.addEventListener('dragover', (e) => {
        e.preventDefault();
        dropOverlay.classList.add('active');
    });

    window.addEventListener('drop', (e) => {
        e.preventDefault();
        dropOverlay.classList.remove('active');
        handleFiles(e.dataTransfer.files);
    });

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
        const row = document.createElement('div');
        row.classList.add('message-row', 'message-sent');

        row.innerHTML = createMessageHTML(content, time, type);

        document.querySelector('.messages-container').appendChild(row);
        document.querySelector('.messages-container').scrollTop = document.querySelector('.messages-container').scrollHeight;
    }

    contactCard.forEach(card => {
        card.addEventListener('click', () => {
            contactCard.forEach(c => c.classList.remove('active'));

            card.classList.add('active');

            const name = card.querySelector('.contact-name').innerText;
            chatHeaderName.innerText = name;

            messagesContainer.innerHTML = "";
        });
    });

    contactSearch,addEventListener('input', (e) => {
        const term = e.target.value.toLowerCase();

        contactCard.forEach(card => {
            const name = card.querySelector('.contact-name').innerText.toLowerCase();
            if (name.includes(term)) {
                card.style.display = "flex";
            } else {
                card.style.display = 'none';
            }
        });
    });

    addContactBtn.onclick = () => addContactModal.style.display = 'flex';
    closeContactModal.onclick = () => addContactModal.style.display = 'none';

    function setupContactClick(card) {
    card.addEventListener('click', () => {
        
        document.querySelectorAll('.contact-card').forEach(c => c.classList.remove('active'));
        
        card.classList.add('active');
        
        const nameElement = card.querySelector('.contact-name');
        if (nameElement && chatHeaderName) {
            chatHeaderName.innerText = nameElement.innerText;
        }
        
        if (messagesContainer) {
            messagesContainer.innerHTML = "";
        }
    });
}

function createContactCardHTML(name, lastMsg = "Nova conversa...") {
    const card = document.createElement('div');
    card.classList.add('contact-card');
    
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
    
    const btn = card.querySelector('.contact-options-btn');
    if (btn) {
        btn.onclick = (e) => {
            e.stopPropagation();
            const menu = card.querySelector('.contact-menu');
            document.querySelectorAll('.contact-menu').forEach(m => {
                if (m !== menu) m.classList.remove('show');
            });
            menu.classList.toggle('show');
        };
    }

    setupContactClick(card);
    
    return card;
}

    saveContactBtn.onclick = () => {
        const name = document.getElementById('new-contact-name').value.trim();
        if (name) {
            const newCard = createContactCardHTML(name);
            contactsList.prepend(newCard); 
            
            document.getElementById('new-contact-name').value = '';
            document.getElementById('new-contact-number').value = '';
            addContactModal.style.display = 'none';
        }
    };

    document.querySelectorAll('.contact-card').forEach(card => {
        card.onclick = () => {
            document.querySelectorAll('.contact-card').forEach(c => c.classList.remove('active'));
            card.classList.add('active');
            const name = card.querySelector('.contact-name').innerText;
            document.querySelector('.chat-header-info strong').innerText = name;
        };
    });

    document.querySelectorAll('.contact-card').forEach(card => {

    if (!card.querySelector('.contact-options-wrapper')) {
        const optionsHTML = `
            <div class="contact-options-wrapper">
                <button class="contact-options-btn" type="button">⋮</button>
                <div class="contact-menu">
                    <button onclick="editContact(this)">Editar Nome</button>
                    <button onclick="clearChat(this)">Limpar Conversa</button>
                    <hr>
                    <button class="delete-contact-btn" onclick="deleteContact(this)">Excluir Contato</button>
                </div>
            </div>`;
        
        card.insertAdjacentHTML('beforeend', optionsHTML);

        const btn = card.querySelector('.contact-options-btn');
        btn.onclick = (e) => {
            e.stopPropagation();
            const menu = card.querySelector('.contact-menu');
            document.querySelectorAll('.contact-menu').forEach(m => {
                if (m !== menu) m.classList.remove('show');
            });
            menu.classList.toggle('show');
        };
    }
});

function handleMobileView() {
    if (window.innerWidth <= 768) {
        const sidebar = document.querySelector('.sidebar');
        sidebar.classList.add('hidden');
    }
}

document.querySelectorAll('.contact-card').forEach(card => {
    card.addEventListener('click', () => {
        handleMobileView();
    });
});

const toggleBtn = document.getElementById('toggle-sidebar');
toggleBtn.addEventListener('click', () => {
    const sidebar = document.querySelector('.sidebar');
    if (window.innerWidth <= 768) {
        sidebar.classList.remove('hidden');
    }
});
});