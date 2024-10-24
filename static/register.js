// 회원가입 함수
function registerUser(event) {
    event.preventDefault(); // 폼이 제출되면 페이지가 새로고침되는 것을 막음

    const userid = document.getElementById('register-userid').value;
    const password = document.getElementById('register-password').value;

    fetch('/api/register', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({
            userid: userid,
            password: password
        })
    })
        .then(response => response.json())
        .then(data => {
            if (data.message === '회원가입 성공') {
                alert('회원가입이 완료되었습니다!');
                window.location.href = '/login';
            } else {
                alert('회원가입 실패: ' + data.message);
            }
        })
        .catch(error => console.error('회원가입 중 오류 발생:', error));
}