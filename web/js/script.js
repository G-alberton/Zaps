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

function handleReply(btn) { console.log("Responder:", btn.closest('.message-bubble').querySelector('.message-text').innerText); }
function handleForward(btn) { console.log("Encaminhar"); }

document.addEventListener('DOMContentLoaded', () => {
    const chatInput = document.querySelector('.chat-input-area textarea') || document.querySelector('.chat-input-area input');
    const sendBtn = document.querySelector('.send-button');
    const messagesContainer = document.querySelector('.messages-container');
    const searchInput = document.getElementById('message-search');
    const searchBtn = document.getElementById('search-button');

    console.log("Sistema de chat iniciado!");

    if (chatInput.tagName.toLowerCase() === 'textarea') {
        chatInput.addEventListener('input', function() {
            this.style.height = 'auto'; 
            this.style.height = (this.scrollHeight) + 'px'; 
        });
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
            chatInput.focus();
            scrollToBottom();

            // Simulação do Bot
            setTimeout(() => {
                receiveMessage("Resposta automática: " + messageText);
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

    function scrollToBottom() {
        messagesContainer.scrollTop = messagesContainer.scrollHeight;
    }

    if (sendBtn) sendBtn.addEventListener('click', sendMessage);

    if (chatInput) {
        chatInput.addEventListener('keydown', (e) => {
            if (e.key === 'Enter' && !e.shiftKey) {
                e.preventDefault();
                sendMessage();
            }
        });
    }

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

    function performSearch() {
        if (!searchInput) return;
    const searchTerm = searchInput.value.toLowerCase().trim();
    const allMessages = document.querySelectorAll('.message-row');

    allMessages.forEach(row => {
        const textElement = row.querySelector('.message-text');
        const bubble = row.querySelector('.message-bubble');
        
        if (textElement && bubble) {
            const messageText = textElement.innerText.toLowerCase();

            if (messageText.includes(searchTerm)) {
                row.style.display = 'flex'; 
                
                if (searchTerm !== "") {
                    bubble.style.backgroundColor = '#fff3cd';
                    bubble.style.color = '#333'; 
                } else {
                    bubble.style.backgroundColor = '';
                    bubble.style.color = ''; 
                }
            } else {
                row.style.display = 'none';
            }
        }
    });
    }

    if (searchBtn) searchBtn.addEventListener('click', performSearch);
    if (searchInput) {
        searchInput.addEventListener('keydown', (e) => { if (e.key === 'Enter') performSearch(); });
        searchInput.addEventListener('input', () => { if (searchInput.value === "") performSearch(); });
    }
});