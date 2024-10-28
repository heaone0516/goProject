document.addEventListener('DOMContentLoaded', function () {
    // 회원가입 함수
    function registerUser(event) {
        event.preventDefault(); // 폼이 제출되면 페이지가 새로고침되는 것을 막음

        const userid = document.getElementById('register-userid').value;
        const password = document.getElementById('register-password').value;
        const confirmPassword = document.getElementById('register-confirm-password').value;

        // 비밀번호와 비밀번호 확인이 일치하는지 최종적으로 확인
        if (password !== confirmPassword) {
            alert('비밀번호가 일치하지 않습니다. 다시 입력해주세요.');
            return;
        }

        // 비밀번호 일치 확인 후 API 호출
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
                if (!response.ok) {
                    return response.json().then(data => {
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

    // 비밀번호 확인 필드의 입력을 실시간으로 체크하여 비밀번호 일치 여부를 표시
    document.getElementById('register-confirm-password').addEventListener('input', function() {
        const password = document.getElementById('register-password').value;
        const confirmPassword = document.getElementById('register-confirm-password').value;
        const message = document.getElementById('password-match-message');

        // 비밀번호가 일치하지 않으면 메시지를 표시
        if (password !== confirmPassword) {
            message.style.display = 'inline';
        } else {
            message.style.display = 'none';
        }
    });

    // 회원가입 폼의 submit 이벤트 리스너 등록
    document.querySelector('form').addEventListener('submit', registerUser);
});