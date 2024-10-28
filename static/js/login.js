// 로그인 후 토큰 저장
function loginUser() {
    const userid = document.getElementById('login-userid').value;
    const password = document.getElementById('login-password').value;

    fetch('/api/login', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ userid, password })
    })
        .then(response => response.json())
        .then(data => {
            if (data.token) {
                // 토큰 저장 (예: localStorage)
                localStorage.setItem('token', data.token);
                alert('로그인 성공!');
                window.location.href = '/';
            } else {
                alert('로그인 실패');
            }
        });
}

// API 호출 시 토큰 포함
function fetchProtectedData() {
    const token = localStorage.getItem('token');
    fetch('/api/protected', {
        method: 'GET',
        headers: {
            'Authorization': `Bearer ${token}`
        }
    })
        .then(response => response.json())
        .then(data => {
            console.log(data);
        });
}