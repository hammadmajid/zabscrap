// DOM Elements
const loginView = document.getElementById('loginView');
const loadingView = document.getElementById('loadingView');
const resultsView = document.getElementById('resultsView');
const loginForm = document.getElementById('loginForm');
const backBtn = document.getElementById('backBtn');
const coursesContainer = document.getElementById('coursesContainer');
const errorToast = document.getElementById('errorToast');
const errorMessage = document.getElementById('errorMessage');

// View Management
function showView(view) {
    loginView.classList.add('hidden');
    loadingView.classList.add('hidden');
    resultsView.classList.add('hidden');
    view.classList.remove('hidden');
}

// Error Handling
function showError(message) {
    errorMessage.textContent = message;
    errorToast.classList.remove('hidden');
    
    setTimeout(() => {
        errorToast.classList.add('hidden');
    }, 5000);
}

// Form Submission
loginForm.addEventListener('submit', async (e) => {
    e.preventDefault();
    
    const username = document.getElementById('username').value;
    const password = document.getElementById('password').value;
    
    if (!username || !password) {
        showError('Please enter both username and password');
        return;
    }
    
    // Show loading state
    showView(loadingView);
    
    try {
        const response = await fetch('/fetch', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ username, password }),
        });
        
        const data = await response.json();
        
        if (!response.ok) {
            throw new Error(data.error || 'Failed to fetch attendance data');
        }
        
        if (data.success) {
            renderCourses(data.data);
            showView(resultsView);
        } else {
            throw new Error(data.error || 'Unknown error occurred');
        }
    } catch (error) {
        console.error('Error:', error);
        showError(error.message);
        showView(loginView);
    }
});

// Back Button
backBtn.addEventListener('click', () => {
    loginForm.reset();
    coursesContainer.innerHTML = '';
    showView(loginView);
});

// Render Courses
function renderCourses(courses) {
    coursesContainer.innerHTML = '';
    
    if (!courses || courses.length === 0) {
        coursesContainer.innerHTML = `
            <div class="card">
                <p style="text-align: center; font-weight: 600; font-size: 1.2rem;">
                    No attendance data found
                </p>
            </div>
        `;
        return;
    }
    
    courses.forEach(course => {
        const courseCard = createCourseCard(course);
        coursesContainer.appendChild(courseCard);
    });
}

// Create Course Card
function createCourseCard(course) {
    const card = document.createElement('div');
    card.className = 'course-card';
    
    const header = document.createElement('div');
    header.className = 'course-header';
    header.innerHTML = `
        <h3 class="course-name">${escapeHtml(course.courseName)}</h3>
        <p class="course-instructor">Instructor: ${escapeHtml(course.instructor)}</p>
    `;
    
    const body = document.createElement('div');
    body.className = 'course-body';
    
    const tableContainer = document.createElement('div');
    tableContainer.className = 'table-container';
    
    const table = createAttendanceTable(course.records);
    tableContainer.appendChild(table);
    
    body.appendChild(tableContainer);
    card.appendChild(header);
    card.appendChild(body);
    
    return card;
}

// Create Attendance Table
function createAttendanceTable(records) {
    const table = document.createElement('table');
    
    // Table Header
    const thead = document.createElement('thead');
    thead.innerHTML = `
        <tr>
            <th>Lecture #</th>
            <th>Date</th>
            <th>Status</th>
        </tr>
    `;
    table.appendChild(thead);
    
    // Table Body
    const tbody = document.createElement('tbody');
    
    if (!records || records.length === 0) {
        tbody.innerHTML = `
            <tr>
                <td colspan="3" style="text-align: center; padding: 20px;">
                    No attendance records
                </td>
            </tr>
        `;
    } else {
        records.forEach(record => {
            const row = document.createElement('tr');
            
            const statusClass = record.status.toLowerCase() === 'absent' 
                ? 'status-absent' 
                : 'status-present';
            
            row.innerHTML = `
                <td><strong>${escapeHtml(record.lecture)}</strong></td>
                <td>${escapeHtml(record.date)}</td>
                <td>
                    <span class="status-badge ${statusClass}">
                        ${escapeHtml(record.status)}
                    </span>
                </td>
            `;
            
            tbody.appendChild(row);
        });
    }
    
    table.appendChild(tbody);
    return table;
}

// Utility: Escape HTML to prevent XSS
function escapeHtml(unsafe) {
    if (!unsafe) return '';
    return unsafe
        .toString()
        .replace(/&/g, "&amp;")
        .replace(/</g, "&lt;")
        .replace(/>/g, "&gt;")
        .replace(/"/g, "&quot;")
        .replace(/'/g, "&#039;");
}

// Initialize
document.addEventListener('DOMContentLoaded', () => {
    showView(loginView);
});
