// 회원가입 함수
function registerUser(event) {
    event.preventDefault(); // 폼 제출 시 페이지 새로고침 방지

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
        .then(response => {
            // JSON 응답을 처리하기 전에 상태 코드를 확인
            if (!response.ok) {
                return response.json().then(data => {
                    // 상태 코드에 따른 오류 메시지 처리
                    if (data.error_code) {
                        switch (data.error_code) {
                            case 'USER_EXISTS':
                                throw new Error('이미 존재하는 사용자 ID입니다.');
                            case 'DB_ERROR':
                                throw new Error('데이터베이스 오류가 발생했습니다.');
                            default:
                                throw new Error('알 수 없는 오류가 발생했습니다.');
                        }
                    } else {
                        throw new Error('오류가 발생했습니다.');
                    }
                });
            }
            return response.json();
        })
        .then(data => {
            if (data.message === '회원가입 성공') {
                alert('회원가입이 완료되었습니다!');
                window.location.href = '/';
            }
        })
        .catch(error => {
            alert('회원가입 실패: ' + error.message);
            console.error('회원가입 중 오류 발생:', error);
        });
}