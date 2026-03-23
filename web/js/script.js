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

function handleReply(btn) { console.log("Responder clicado"); }

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

    let timerInterval;
    let seconds = 0;

    function createMessageHTML(content, time, type = 'text') {
        const isAudio = type === 'audio';
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
                        <button class="audio-play-btn">▶</button>
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
    });

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
        document.querySelectorAll('.message-row').forEach(row => {
            const text = row.querySelector('.message-text')?.innerText.toLowerCase() || "";
            row.style.display = text.includes(term) || term === "" ? 'flex' : 'none';
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
});