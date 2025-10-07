class TimerController {
    constructor() {
        this.secret = localStorage.getItem('timerSecret') || '';
        this.baseURL = window.location.origin;
        this.initElements();
        this.bindEvents();
        this.updateSecretDisplay();
        
        // 自动刷新状态
        this.startAutoRefresh();
    }

    initElements() {
        this.secretInput = document.getElementById('secret');
        this.saveSecretBtn = document.getElementById('save-secret');
        this.refreshStatusBtn = document.getElementById('refresh-status');
        this.sendSignalBtn = document.getElementById('send-signal');
        this.remainingTimeEl = document.getElementById('remaining-time');
        this.progressFillEl = document.getElementById('progress-fill');
        this.notificationEl = document.getElementById('notification');
    }

    bindEvents() {
        this.saveSecretBtn.addEventListener('click', () => this.saveSecret());
        this.refreshStatusBtn.addEventListener('click', () => this.fetchStatus());
        this.sendSignalBtn.addEventListener('click', () => this.sendSignal());
        
        // 回车键保存密钥
        this.secretInput.addEventListener('keypress', (e) => {
            if (e.key === 'Enter') {
                this.saveSecret();
            }
        });
    }

    saveSecret() {
        const secret = this.secretInput.value.trim();
        if (!secret) {
            this.showNotification('请输入密钥', 'error');
            return;
        }
        
        this.secret = secret;
        localStorage.setItem('timerSecret', secret);
        this.updateSecretDisplay();
        this.showNotification('密钥已保存', 'success');
        this.fetchStatus(); // 保存后自动刷新状态
    }

    updateSecretDisplay() {
        if (this.secret) {
            this.secretInput.placeholder = '密钥已保存';
            this.secretInput.value = '';
        } else {
            this.secretInput.placeholder = '输入您的访问密钥';
        }
    }

    async fetchStatus() {
        try {
            const response = await fetch(`${this.baseURL}/timer/status`);
            const data = await response.json();
            
            if (data.code === 200) {
                this.displayTime(data.remaining_time);
            } else {
                this.showNotification('获取状态失败', 'error');
            }
        } catch (error) {
            console.error('Error fetching status:', error);
            this.showNotification('网络错误，请稍后重试', 'error');
        }
    }

    displayTime(seconds) {
        const hours = Math.floor(seconds / 3600);
        const minutes = Math.floor((seconds % 3600) / 60);
        const secs = Math.floor(seconds % 60);
        
        this.remainingTimeEl.textContent = 
            `${hours.toString().padStart(2, '0')}:${minutes.toString().padStart(2, '0')}:${secs.toString().padStart(2, '0')}`;
        
        // 更新进度条（假设最大时间为24小时）
        const maxTime = 24 * 3600; // 24小时
        const progress = Math.min(100, (seconds / maxTime) * 100);
        this.progressFillEl.style.width = `${100 - progress}%`;
    }

    async sendSignal() {
        if (!this.secret) {
            this.showNotification('请先设置密钥', 'error');
            return;
        }
        
        try {
            const response = await fetch(`${this.baseURL}/signal?secret=${encodeURIComponent(this.secret)}`);
            const data = await response.json();
            
            if (data.code === 200) {
                this.showNotification('信号发送成功！', 'success');
                setTimeout(() => this.fetchStatus(), 1000); // 1秒后刷新状态
            } else {
                this.showNotification(`操作失败: ${data.message}`, 'error');
            }
        } catch (error) {
            console.error('Error sending signal:', error);
            this.showNotification('网络错误，请稍后重试', 'error');
        }
    }

    showNotification(message, type) {
        this.notificationEl.textContent = message;
        this.notificationEl.className = `notification ${type} show`;
        
        setTimeout(() => {
            this.notificationEl.classList.remove('show');
        }, 3000);
    }

    startAutoRefresh() {
        // 每30秒自动刷新一次状态
        setInterval(() => {
            if (this.secret) {
                this.fetchStatus();
            }
        }, 30000);
        
        // 页面加载完成后立即获取一次状态
        setTimeout(() => this.fetchStatus(), 1000);
    }
}

// 初始化应用
document.addEventListener('DOMContentLoaded', () => {
    new TimerController();
});