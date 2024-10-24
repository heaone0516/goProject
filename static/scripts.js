// DOMContentLoaded
document.addEventListener('DOMContentLoaded', function () {
    // 검색 폼 제출 시 검색어로 게시물 목록을 필터링
    const searchForm = document.getElementById('search-form');
    if (searchForm) {
        searchForm.addEventListener('submit', function (event) {
            event.preventDefault();
            const searchQuery = document.getElementById('search').value;
            loadPosts(1, searchQuery); // 검색 시 첫 페이지부터 다시 로드
        });
    }

    // 페이지 로드 시 게시물 목록을 불러옴
    loadPosts();

    // 페이지 로드 시 로그인 상태 확인
    checkLoginStatus();

    // 처음 페이지 로드 시 기본 1페이지를 로드
    loadPosts(1); // 기본 1페이지 로드

    // 로그인한 사용자 정보를 받아오는 함수
    getUserInfo();

});


// 게시물 목록을 가져와서 테이블에 렌더링하는 함수
function loadPosts(page = 1, search = "") {
    let url = `/api/posts?page=${page}&limit=10`;

    // 검색어가 있을 경우 URL에 쿼리 파라미터 추가
    if (search) {
        url += `&search=${encodeURIComponent(search)}`;
    }

    // 데이터를 서버에서 fetch로 가져오기
    fetch(url)
        .then(response => {
            if (!response.ok) {
                throw new Error('게시물을 불러오는 데 실패했습니다.');
            }
            return response.json();
        })
        .then(data => {
            console.log(data);  // 응답 전체를 출력해서 확인 (필요한 경우 제거 가능)

            const posts = data.posts;  // posts 목록만 분리
            const totalPages = data.totalPages;  // totalPages 값을 분리

            const tableBody = document.getElementById('posts-table-body');
            tableBody.innerHTML = ""; // 기존 내용을 초기화

            if (!posts || posts.length === 0) {  // 검색 결과가 없을 경우 처리
                tableBody.innerHTML = "<tr><td colspan='6'>검색 결과가 없습니다.</td></tr>";
                return;
            }

            posts.forEach(post => {
                const row = document.createElement('tr');
                row.innerHTML = `
                <td>${post.id}</td>
                <td>${post.title}</td>
                <td>${post.content}</td>
                <td>${post.author}</td>
                <td>${post.created_at}</td>
                <td>
                    <button onclick="editPost(${post.id})">수정</button>
                    <button onclick="deletePost(${post.id})">삭제</button>
                </td>
            `;
                tableBody.appendChild(row);
            });

            // 페이지네이션 생성
            createPagination(totalPages, page, search);
        })
        .catch(error => {
            console.error('게시물을 불러오는 중 오류 발생:', error);
        });
}

// 페이지네이션 버튼을 생성하는 함수
function createPagination(totalPages, currentPage, search) {
    console.log("Total Pages:", totalPages);  // totalPages 값 확인
    const paginationDiv = document.querySelector('.pagination');
    paginationDiv.innerHTML = ""; // 기존 페이지네이션 초기화

    // 페이지 수가 1보다 클 때만 페이지네이션을 표시
    if (totalPages > 1) {
        // 이전 버튼 생성
        const prevButton = document.createElement('a');
        prevButton.innerHTML = '&laquo;';
        prevButton.href = 'javascript:void(0);';
        prevButton.onclick = () => loadPosts(currentPage - 1, search);
        prevButton.style.pointerEvents = currentPage === 1 ? 'none' : 'auto'; // 첫 페이지에서 비활성화
        paginationDiv.appendChild(prevButton);

        // 각 페이지 번호 생성
        for (let i = 1; i <= totalPages; i++) {
            const pageButton = document.createElement('a');
            pageButton.innerHTML = i;
            pageButton.href = 'javascript:void(0);';

            // 현재 페이지인 경우에만 'active' 클래스를 추가
            if (currentPage === i) {
                pageButton.classList.add('active');
            }

            pageButton.onclick = () => loadPosts(i, search);
            paginationDiv.appendChild(pageButton);
        }

        // 다음 버튼 생성
        const nextButton = document.createElement('a');
        nextButton.innerHTML = '&raquo;';
        nextButton.href = 'javascript:void(0);';
        nextButton.onclick = () => loadPosts(currentPage + 1, search);
        nextButton.style.pointerEvents = currentPage === totalPages ? 'none' : 'auto'; // 마지막 페이지에서 비활성화
        paginationDiv.appendChild(nextButton);
    }
}


// 글 작성 폼을 보여주는 함수
function createPostForm() {
    // 폼 필드를 초기화
    document.getElementById('create-title').value = '';   // 제목 필드 초기화
    document.getElementById('create-content').value = ''; // 내용 필드 초기화

    // 폼을 보여줍니다.
    document.getElementById('create-post-form').style.display = 'block';
}

// 게시물 생성 함수
function createPost() {
    const title = document.getElementById('create-title').value;
    const content = document.getElementById('create-content').value;
    const author = document.getElementById('create-author').value;

    const postData = {
        title: title,
        content: content,
        author: author
    };

    fetch('/api/posts', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(postData)
    })
        .then(response => {
            if (response.status === 201) {
                alert('게시물이 성공적으로 생성되었습니다.');
                loadPosts(); // 게시물 목록을 갱신합니다.
                hidePostForm('create'); // 폼을 숨깁니다.
            } else {
                throw new Error('게시물 생성에 실패했습니다.');
            }
        })
        .catch(error => {
            console.error('게시물 생성 중 오류 발생:', error);
        });
}

// 게시물 수정 폼에 데이터 채워 넣기
function editPost(postId) {
    fetch(`/api/posts/edit?id=${postId}`)
        .then(response => {
            if (!response.ok) {
                throw new Error('게시물을 불러오는 데 실패했습니다.');
            }
            return response.json();
        })
        .then(post => {
            // 폼을 열고 기존 게시물 데이터를 폼에 채워 넣음
            document.getElementById('form-title').innerText = '게시글 수정'; // 제목 변경
            document.getElementById('post-id').value = post.id; // 게시물 ID 저장
            document.getElementById('edit-title').value = post.title; // 기존 제목
            document.getElementById('edit-content').value = post.content; // 기존 내용
            document.getElementById('edit-author').value = post.author; // 기존 작성자

            showPostForm('edit'); // 수정 폼을 표시
        })
        .catch(error => {
            console.error('게시물을 불러오는 중 오류 발생:', error);
        });
}

// 게시물 업데이트 함수
function updatePost(postId) {
    const title = document.getElementById('edit-title').value;
    const content = document.getElementById('edit-content').value;

    const postData = {
        title: title,
        content: content
    };

    fetch(`/api/posts/update?id=${postId}`, {
        method: 'PUT',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(postData)
    })
        .then(response => {
            if (response.ok) {
                alert('게시물이 성공적으로 수정되었습니다.');
                loadPosts(); // 게시물 목록 갱신
                hidePostForm('edit'); // 폼 숨기기
                clearPostForm('edit'); // 수정 폼 초기화
            } else {
                throw new Error('게시물 수정에 실패했습니다.');
            }
        })
        .catch(error => {
            console.error('게시물 수정 중 오류 발생:', error);
        });
}

// 게시물 생성 또는 수정 버튼을 통해 요청 보냄
function createOrUpdatePost() {
    const postId = document.getElementById('post-id').value;
    if (postId) {
        updatePost(postId); // postId가 있으면 수정
    } else {
        createPost(); // 없으면 새 게시물 작성
    }
}

// 게시물 작성/수정 폼을 표시하는 함수
function showPostForm() {
    document.getElementById('edit-post-form').style.display = 'block';
}

// 게시물 작성/수정 폼을 숨기는 함수
function hidePostForm(formType) {
    if (formType === 'create') {
        document.getElementById('create-post-form').style.display = 'none'; // 새 글 작성 폼을 숨깁니다.
    } else if (formType === 'edit') {
        document.getElementById('edit-post-form').style.display = 'none'; // 게시물 수정 폼을 숨깁니다.
    }
}


// 게시물 삭제 함수
function deletePost(postId) {
    if (!confirm('정말 삭제하시겠습니까?')) {
        return;
    }

    fetch(`/api/posts/delete?id=${postId}`, {
        method: 'DELETE'
    })
        .then(response => {
            if (response.ok) {
                alert('게시물이 성공적으로 삭제되었습니다.');
                loadPosts(); // 게시물 목록 갱신
            } else {
                throw new Error('게시물 삭제에 실패했습니다.');
            }
        })
        .catch(error => {
            console.error('게시물 삭제 중 오류 발생:', error);
        });
}


// 게시물 작성 폼 초기화 함수
function clearPostForm(formType) {
    if (formType === 'create') {
        document.getElementById('create-title').value = '';
        document.getElementById('create-content').value = '';
        document.getElementById('create-author').value = '';
    } else if (formType === 'edit') {
        document.getElementById('edit-title').value = '';
        document.getElementById('edit-content').value = '';
        document.getElementById('edit-author').value = '';
        document.getElementById('post-id').value = '';
        document.getElementById('edit-submit-button').innerText = '글 작성';
    }
}


// 로그아웃 함수 정의
function logoutUser() {
    fetch('/api/logout', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        }
    })
        .then(response => response.json())
        .then(data => {
            if (data.message === '로그아웃 성공') {
                alert('로그아웃이 완료되었습니다.');
                window.location.href = '/index'; // 로그아웃 후 로그인 페이지로 리다이렉트
            } else {
                alert('로그아웃 실패: ' + data.message);
            }
        })
        .catch(error => console.error('로그아웃 중 오류 발생:', error));
}

// 로그인 여부 확인 함수
function checkLoginStatus() {
    fetch('/api/is_logged_in', {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json'
        }
    })
        .then(response => response.json())
        .then(data => {
            if (data.is_logged_in) {
                // 로그인 상태라면 로그아웃 버튼 표시
                document.getElementById('login-section').style.display = 'none';
                document.getElementById('logout-section').style.display = 'block';
            } else {
                // 로그인 상태가 아니면 로그인/회원가입 버튼 표시
                document.getElementById('login-section').style.display = 'block';
                document.getElementById('logout-section').style.display = 'none';
            }
        })
        .catch(error => console.error('로그인 상태 확인 중 오류 발생:', error));
}

// 로그인한 사용자 정보를 받아오는 함수
function getUserInfo() {
    fetch('/api/get_user')
        .then(response => response.json())
        .then(data => {
            if (data.user_id) {
                // 작성자 필드에 사용자 아이디 자동 입력
                document.getElementById('user-id').innerText = data.user_id;
                document.getElementById('create-author').value = data.user_id;
            }
        })
        .catch(error => console.error('사용자 정보를 불러오는 중 오류 발생:', error));
}


