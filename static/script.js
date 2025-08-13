document.addEventListener('DOMContentLoaded', () => {
    const projectHeroSection = document.getElementById('project-hero-section');
    const projectMainTitle = document.getElementById('project-main-title');
    const projectStoryDescription = document.getElementById('project-story-description');
    const projectMediaDisplay = document.getElementById('project-media-display');
    const beforeBtn = document.getElementById('before-btn');
    const afterBtn = document.getElementById('after-btn');
    const futureStateLink = document.getElementById('future-state-link');
    const projectCost = document.getElementById('project-cost');
    const projectFunded = document.getElementById('project-funded');
    const projectFundedOfCost = document.getElementById('project-funded-of-cost');
    const projectProgressPercentage = document.getElementById('project-progress-percentage');
    const projectProgressBar = document.getElementById('project-progress-bar');
    const ongoingProjectsContainer = document.getElementById('ongoing-projects-container');
    const completedProjectsContainer = document.getElementById('completed-projects-container');

    let allProjects = [];
    let currentProject = null;

    // Function to fetch project data
    async function fetchProjects() {
        try {
            const response = await fetch('project.json');
            const data = await response.json();
            allProjects = data.projects;
            displayProjects(allProjects);
            if (allProjects.length > 0) {
                loadProjectDetails(allProjects[0]); // Load the first project by default
            }
        } catch (error) {
            console.error('Error fetching projects:', error);
        }
    }

    // Function to display project cards
    function displayProjects(projects) {
        ongoingProjectsContainer.innerHTML = '';
        completedProjectsContainer.innerHTML = '';

        projects.forEach(project => {
            const cardHtml = `
                <div class="col-12 col-md-6">
                    <div class="card h-100 project-card" data-project-id="${project.id}">
                        <img src="${project.before.image}" class="card-img-top" alt="${project.title}">
                        <div class="card-body d-flex flex-column">
                            <h2 class="card-title h5 fw-bold">${project.title}</h2>
                            <p class="card-text text-small">${project.description.substring(0, 100)}.....</p>
                            <p class="card-text small mb-2"><strong>${project.date}</strong></p>
                            ${project.status === 'ongoing' ? `
                            <div class="mt-auto">
                                <div class="progress mb-2">
                                    <div class="progress-bar" role="progressbar" style="width: ${project.progress}%;" aria-valuenow="${project.progress}" aria-valuemin="0" aria-valuemax="100"></div>
                                </div>
                                <p class="card-text fw-bold text-start">N${project.funded} out of N${project.cost}</p>
                            </div>
                            ` : ''}
                        </div>
                    </div>
                </div>
            `;
            if (project.status === 'ongoing') {
                ongoingProjectsContainer.innerHTML += cardHtml;
            } else {
                completedProjectsContainer.innerHTML += cardHtml;
            }
        });

        // Add event listeners to project cards
        document.querySelectorAll('.project-card').forEach(card => {
            card.addEventListener('click', (event) => {
                const projectId = event.currentTarget.dataset.projectId;
                const selectedProject = allProjects.find(p => p.id == projectId);
                if (selectedProject) {
                    loadProjectDetails(selectedProject);
                    // Scroll to the top of the page or project hero section
                    projectHeroSection.scrollIntoView({ behavior: 'smooth' });
                }
            });
        });
    }

    // Function to load project details into the main display area
    function loadProjectDetails(project) {
        currentProject = project;
        projectMainTitle.textContent = project.title;
        projectStoryDescription.textContent = project.description; // Assuming the story is the full description
        projectCost.textContent = `N${project.cost}`;
        projectFunded.textContent = `N${project.funded}`;
        projectFundedOfCost.textContent = `of N${project.cost} Funded`;
        projectProgressPercentage.textContent = `${project.progress}%`;
        projectProgressBar.style.width = `${project.progress}%`;
        projectProgressBar.setAttribute('aria-valuenow', project.progress);
        futureStateLink.href = project.future_state_3d;

        // Set initial media to 'Before'
        displayBeforeMedia(project);
        beforeBtn.classList.add('active');
        afterBtn.classList.remove('active');
    }

    // Function to display 'Before' media
    function displayBeforeMedia(project) {
        projectMediaDisplay.innerHTML = ''; // Clear previous content
        let mediaHtml = '';
        if (project.before.image) {
            mediaHtml += `
                <div class="media-item mb-3">
                    <img src="${project.before.image}" alt="Project Before Image" class="img-fluid rounded shadow">
                </div>
            `;
        }
        if (project.before.video) {
            mediaHtml += `
                <div class="media-item mb-3">
                    <div class="embed-responsive embed-responsive-16by9 rounded shadow">
                        <iframe class="embed-responsive-item" src="${project.before.video}" allowfullscreen></iframe>
                    </div>
                </div>
            `;
        }
        projectMediaDisplay.innerHTML = mediaHtml;
        projectHeroSection.style.setProperty('--page-background-image', `url('${project.before.background_image}')`);
        futureStateLink.textContent = "Videos Of Project";
    }

    // Function to display 'After' media
    function displayAfterMedia(project) {
        projectMediaDisplay.innerHTML = ''; // Clear previous content
        let mediaHtml = '';
        if (project.after.images && project.after.images.length > 0) {
            project.after.images.forEach(image => {
                mediaHtml += `
                    <div class="media-item mb-3">
                        <img src="${image}" alt="Project After Image" class="img-fluid rounded">
                    </div>
                `;
            });
        }
        projectMediaDisplay.innerHTML = mediaHtml;
        projectHeroSection.style.setProperty('--page-background-image', `url('${project.after.background_image}')`);
        futureStateLink.textContent = "Future State 3D View";
    }

    // Event listeners for Before/After buttons
    beforeBtn.addEventListener('click', (e) => {
        e.preventDefault();
        if (currentProject) {
            displayBeforeMedia(currentProject);
            beforeBtn.classList.add('active');
            afterBtn.classList.remove('active');
        }
    });

    afterBtn.addEventListener('click', (e) => {
        e.preventDefault();
        if (currentProject) {
            displayAfterMedia(currentProject);
            afterBtn.classList.add('active');
            beforeBtn.classList.remove('active');
        }
    });

    // Initial fetch of projects
    fetchProjects();
});