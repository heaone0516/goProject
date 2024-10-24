function loginUser() {
    const userid = document.getElementById('login-userid').value;
    const password = document.getElementById('login-password').value;

    fetch('/api/login', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'  // JSON 형식으로 전송
        },
        body: JSON.stringify({
            userid: userid,
            password: password
        })
    })
        .then(response => {
            if (!response.ok) {
                return response.json().then(err => { throw new Error(err.message); });
            }
            return response.json();
        })
        .then(data => {
            if (data.message === '로그인 성공') {
                alert('로그인 성공!');
                window.location.href = '/'; // 메인 페이지로 리다이렉트
            } else {
                alert('로그인 실패: ' + data.message);
            }
        })
        .catch(error => console.error('로그인 중 오류 발생:', error));
}