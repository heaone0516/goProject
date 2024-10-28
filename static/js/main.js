// 게시물 목록을 가져와서 테이블에 렌더링하는 함수
function loadPosts(page = 1, search = "") {
    let url = `/api/posts?page=${page}&limit=10`;

    if (search) {
        url += `&search=${encodeURIComponent(search)}`;
    }

    fetch(url)
        .then(response => {
            if (!response.ok) {
                throw new Error('게시물을 불러오는 데 실패했습니다.');
            }
            return response.json();
        })
        .then(data => {
            const posts = data.posts;
            const totalPages = data.totalPages;

            const tableBody = document.getElementById('posts-table-body');
            tableBody.innerHTML = "";

            if (!posts || posts.length === 0) {
                tableBody.innerHTML = "<tr><td colspan='6'>검색 결과가 없습니다.</td></tr>";
                return;
            }

            posts.forEach(post => {
                const row = document.createElement('tr');
                row.innerHTML = `
                <td>${post.id}</td>
                <td><a href="javascript:void(0);" onclick="fetchPostDetails(${post.id})">${post.title}</a></td>
                <td>${post.author}</td>
                <td>${post.created_at}</td>
                <td>
                    <button onclick="editPost(${post.id})">수정</button>
                    <button onclick="deletePost(${post.id})">삭제</button>
                </td>`;
                tableBody.appendChild(row);
            });

            createPagination(totalPages, page, search);
        })
        .catch(error => {
            console.error('게시물을 불러오는 중 오류 발생:', error);
        });
}

// 특정 게시물 상세 내용을 불러오는 함수
function fetchPostDetails(postId) {
    fetch(`/api/posts/${postId}`)
        .then(response => {
            if (!response.ok) {
                throw new Error('게시물 내용을 불러오는 데 실패했습니다.');
            }
            return response.json();
        })
        .then(post => {
            const postDetails = document.getElementById('post-details');
            postDetails.innerHTML = `
                <h2 id="post-title">${post.title}</h2>
                <p id="post-content">${post.content}</p>`;
            postDetails.style.display = 'block';
        })
        .catch(error => {
            console.error('게시물 내용을 불러오는 중 오류 발생:', error);
        });
}

// 페이지네이션 버튼을 생성하는 함수
function createPagination(totalPages, currentPage, search) {
    const paginationDiv = document.querySelector('.pagination');
    paginationDiv.innerHTML = "";

    if (totalPages > 1) {
        const prevButton = document.createElement('a');
        prevButton.innerHTML = '&laquo;';
        prevButton.href = 'javascript:void(0);';
        prevButton.onclick = () => {
            loadPosts(currentPage - 1, search);
            document.getElementById('post-details').style.display = 'none';
        };
        prevButton.style.pointerEvents = currentPage === 1 ? 'none' : 'auto';
        paginationDiv.appendChild(prevButton);

        for (let i = 1; i <= totalPages; i++) {
            const pageButton = document.createElement('a');
            pageButton.innerHTML = i;
            pageButton.href = 'javascript:void(0);';
            if (currentPage === i) {
                pageButton.classList.add('active');
            }
            pageButton.onclick = () => {
                loadPosts(i, search);
                document.getElementById('post-details').style.display = 'none';
            };
            paginationDiv.appendChild(pageButton);
        }

        const nextButton = document.createElement('a');
        nextButton.innerHTML = '&raquo;';
        nextButton.href = 'javascript:void(0);';
        nextButton.onclick = () => {
            loadPosts(currentPage + 1, search);
            document.getElementById('post-details').style.display = 'none';
        };
        nextButton.style.pointerEvents = currentPage === totalPages ? 'none' : 'auto';
        paginationDiv.appendChild(nextButton);
    }
}

// DOMContentLoaded for Main Page
document.addEventListener('DOMContentLoaded', function () {
    const searchForm = document.getElementById('search-form');
    if (searchForm) {
        searchForm.addEventListener('submit', function (event) {
            event.preventDefault();
            const searchQuery = document.getElementById('search').value;
            loadPosts(1, searchQuery);
            document.getElementById('post-details').style.display = 'none';
        });
    }

    loadPosts(1); // 기본 1페이지 로드
});