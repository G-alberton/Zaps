document.addEventListener('DOMContentLoaded', () => {
    const sendBtn = document.querySelector('.send-button');
    const chatInput = document.querySelector('.chat-input-area input');
    const messagesContainer = document.querySelector('.messages-container');

    console.log("Sistema de chat iniciado!"); 

    function sendMessage() {
        const messageText = chatInput.value.trim();

        if (messageText !== "") {
            const messageRow = document.createElement('div');
            messageRow.classList.add('message-row', 'message-sent');

            const now = new Date();
            const time = now.getHours() + ":" + now.getMinutes().toString().padStart(2, '0');

            messageRow.innerHTML = `
                <div class="message-bubble">
                    ${messageText}
                    <span class="message-time">${time}</span>
                </div>
            `;

            messagesContainer.appendChild(messageRow);

            setTimeout(() => {
            receiveMessage("Recebi sua mensagem: " + messageText);
            }, 2000);

            chatInput.value = "";
            chatInput.focus();
            
            messagesContainer.scrollTop = messagesContainer.scrollHeight;
        }
    }

    if(sendBtn) {
        sendBtn.addEventListener('click', sendMessage);
    }

    if(chatInput) {
        chatInput.addEventListener('keydown', (e) => {
            if (e.key === 'Enter') {
                e.preventDefault(); 
                sendMessage();
            }
        });
    }

    function receiveMessage(text) {
        const messageRow = document.createElement('div');
        messageRow.classList.add('message-row', 'message-received');

        const now = new Date();
        const time = now.getHours() + ":" + now.getMinutes().toString().padStart(2, '0');

        messageRow.innerHTML = `
            <div class="message-bubble">
                ${text}
                <span class="message-time">${time}</span>
        `;

        messagesContainer.appendChild(messageRow);
        messagesContainer.scrollTop = messagesContainer.scrollHeight;
    }

});