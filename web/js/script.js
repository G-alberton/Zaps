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
            alert("Copiado: " + textToCopy);
            bubble.querySelector('.message-menu').classList.remove('show');
        });
    }
}

function handleReply(btn) { console.log("Responder clicado"); }
function handleForward(btn) { console.log("Encaminhar clicado"); }

document.addEventListener('DOMContentLoaded', () => {
    const chatInput = document.getElementById('chat-input');
    const sendBtn = document.querySelector('.send-button');
    const audioBtn = document.querySelector('.audio-button');
    const messagesContainer = document.querySelector('.messages-container');
    const searchInput = document.getElementById('message-search');
    const searchBtn = document.getElementById('search-button');
    const recordingStatus = document.getElementById('recording-status');
    const timerDisplay = document.getElementById('recording-timer');

    let timerInterval;
    let seconds = 0;

    console.log("Sistema de chat iniciado!");

    function sendMessage() {
        const messageText = chatInput.value.trim();
        if (messageText !== "") {
            const messageRow = document.createElement('div');
            messageRow.classList.add('message-row', 'message-sent');

            const now = new Date();
            const time = now.getHours() + ":" + now.getMinutes().toString().padStart(2, '0');

            messageRow.innerHTML = createMessageHTML(messageText, time);
            messagesContainer.appendChild(messageRow);

            chatInput.value = "";
            chatInput.style.height = 'auto';
            toggleInputButtons("");
            chatInput.focus();
            scrollToBottom();

            // Simulação de Resposta
            setTimeout(() => {
                receiveMessage("Resposta automática para: " + messageText);
            }, 2000);
        }
    }

    function receiveMessage(text) {
        const messageRow = document.createElement('div');
        messageRow.classList.add('message-row', 'message-received');
        const now = new Date();
        const time = now.getHours() + ":" + now.getMinutes().toString().padStart(2, '0');
        messageRow.innerHTML = createMessageHTML(text, time);
        messagesContainer.appendChild(messageRow);
        scrollToBottom();
    }

    function createMessageHTML(text, time) {
        return `
            <div class="message-bubble">
                <div class="message-options-btn"><span>&#9013;</span></div>
                <div class="message-menu">
                    <button onclick="handleReply(this)">Responder</button>
                    <button onclick="handleCopy(this)">Copiar</button>
                    <button onclick="handleForward(this)">Encaminhar</button>
                    <hr>
                    <button class="delete-btn" onclick="handleDelete(this)">Apagar</button>
                </div>
                <div class="message-text">${text}</div>
                <span class="message-time">${time}</span>
            </div>
        `;
    }

    function toggleInputButtons(text) {
        if (text.length > 0) {
            audioBtn.style.display = 'none';
            sendBtn.style.display = 'flex';
        } else {
            audioBtn.style.display = 'flex';
            sendBtn.style.display = 'none';
        }
    }

    chatInput.addEventListener('input', function() {
        // Expandir altura
        this.style.height = 'auto';
        this.style.height = (this.scrollHeight) + 'px';
        // Trocar botão
        toggleInputButtons(this.value.trim());
    });

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

    function stopTimer() {
        clearInterval(timerInterval);
    }

    audioBtn.addEventListener('mousedown', (e) => {
        e.preventDefault();
        audioBtn.classList.add('recording');
        audioBtn.innerHTML = "🛑";
        
        chatInput.style.display = 'none';
        recordingStatus.setAttribute('style', 'display: flex !important');
        
        startTimer();
    });

    audioBtn.addEventListener('mouseup', () => {
        audioBtn.classList.remove('recording');
        audioBtn.innerHTML = "🎙️";
        
        chatInput.style.display = 'block';
        recordingStatus.setAttribute('style', 'display: none !important');
        
        const finalTime = timerDisplay.innerText;
        stopTimer();
        
        if (seconds > 0) {
            receiveMessage("🎤 Áudio enviado (" + finalTime + ")");
        }
    });

    function performSearch() {
        const searchTerm = searchInput.value.toLowerCase().trim();
        const allMessages = document.querySelectorAll('.message-row');

        allMessages.forEach(row => {
            const textElement = row.querySelector('.message-text');
            const bubble = row.querySelector('.message-bubble');
            
            if (textElement && bubble) {
                const messageText = textElement.innerText.toLowerCase();
                if (messageText.includes(searchTerm)) {
                    row.style.display = 'flex'; 
                    bubble.style.backgroundColor = (searchTerm !== "") ? '#fff3cd' : '';
                    bubble.style.color = (searchTerm !== "") ? '#333' : '';
                } else {
                    row.style.display = 'none';
                }
            }
        });
    }

    searchBtn.addEventListener('click', performSearch);
    searchInput.addEventListener('keydown', (e) => { if (e.key === 'Enter') performSearch(); });
    searchInput.addEventListener('input', () => { if (searchInput.value === "") performSearch(); });

    function scrollToBottom() {
        messagesContainer.scrollTop = messagesContainer.scrollHeight;
    }

    sendBtn.addEventListener('click', sendMessage);
    chatInput.addEventListener('keydown', (e) => {
        if (e.key === 'Enter' && !e.shiftKey) {
            e.preventDefault();
            sendMessage();
        }
    });

    document.addEventListener('click', (e) => {
        const btn = e.target.closest('.message-options-btn');
        if (btn) {
            const menu = btn.nextElementSibling; 
            menu.classList.toggle('show');
            document.querySelectorAll('.message-menu').forEach(m => {
                if (m !== menu) m.classList.remove('show');
            });
        } else if (!e.target.closest('.message-menu')) {
            document.querySelectorAll('.message-menu').forEach(m => m.classList.remove('show'));
        }
    });
});