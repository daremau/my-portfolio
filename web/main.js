'use strict';

// Theme toggle ---------------------------------------------------------------
const themeToggle = document.getElementById('theme-toggle');

themeToggle?.addEventListener('click', () => {
	const isDark = document.documentElement.classList.toggle('dark');
	localStorage.setItem('theme', isDark ? 'dark' : 'light');
});

// Project filtering ----------------------------------------------------------
const filterButtons = document.querySelectorAll('.filter-btn');
const projectItems = document.querySelectorAll('.project-item');

filterButtons.forEach((button) => {
	button.addEventListener('click', () => {
		const filter = button.getAttribute('data-filter');

		filterButtons.forEach((btn) => {
			btn.classList.remove('filter-btn-active');
			btn.setAttribute('aria-selected', 'false');
		});
		button.classList.add('filter-btn-active');
		button.setAttribute('aria-selected', 'true');

		projectItems.forEach((item) => {
			const category = item.getAttribute('data-category');
			if (filter === 'All' || category === filter) {
				item.style.display = 'block';
				// Re-trigger the fade-in animation.
				item.classList.remove('visible');
				setTimeout(() => item.classList.add('visible'), 10);
			} else {
				item.style.display = 'none';
			}
		});
	});
});

// Scroll animations ----------------------------------------------------------
const observer = new IntersectionObserver(
	(entries) => {
		entries.forEach((entry) => {
			if (entry.isIntersecting) {
				entry.target.classList.add('visible');
			}
		});
	},
	{ threshold: 0.1, rootMargin: '0px 0px -50px 0px' }
);

document.querySelectorAll('.fade-in').forEach((el) => observer.observe(el));
