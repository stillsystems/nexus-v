document.addEventListener('DOMContentLoaded', async () => {
    const templateSelect = document.getElementById('template');
    const generateBtn = document.getElementById('generate-btn');
    const statusPanel = document.getElementById('status-panel');
    const statusText = document.getElementById('status-text');
    const featuresSection = document.getElementById('features-section');
    const featuresList = document.getElementById('features-list');

    let templatesData = [];

    const nameInput = document.getElementById('name');
    const idInput = document.getElementById('identifier');

    function slugify(text) {
        return text.toString().toLowerCase()
            .replace(/\s+/g, '-')
            .replace(/[^\w\-]+/g, '')
            .replace(/\-\-+/g, '-')
            .replace(/^-+/, '')
            .replace(/-+$/, '');
    }

    nameInput.addEventListener('input', () => {
        if (!idInput.dataset.touched) {
            idInput.value = slugify(nameInput.value);
        }
    });

    idInput.addEventListener('input', () => {
        idInput.dataset.touched = "true";
    });

    // Fetch templates on load
    try {
        const response = await fetch('/api/templates');
        const data = await response.json();
        templatesData = data.templates;
        
        templateSelect.innerHTML = '';
        templatesData.forEach(t => {
            const option = document.createElement('option');
            option.value = t.id;
            option.textContent = t.name;
            templateSelect.appendChild(option);
        });

        // Initialize features for first template
        if (templatesData.length > 0) {
            renderFeatures(templatesData[0]);
        }
    } catch (err) {
        console.error('Failed to load templates:', err);
    }

    templateSelect.addEventListener('change', () => {
        const selected = templatesData.find(t => t.id === templateSelect.value);
        renderFeatures(selected);
    });

    function renderFeatures(template) {
        featuresList.innerHTML = '';
        if (!template || !template.features || template.features.length === 0) {
            featuresSection.classList.add('hidden');
            return;
        }

        featuresSection.classList.remove('hidden');
        template.features.forEach(f => {
            const div = document.createElement('div');
            div.className = 'feature-item';
            div.innerHTML = `
                <input type="checkbox" id="feat-${f.id}" ${f.default ? 'checked' : ''} data-feature-id="${f.id}">
                <span>${f.prompt}</span>
            `;
            div.addEventListener('click', (e) => {
                if (e.target.tagName !== 'INPUT') {
                    const cb = div.querySelector('input');
                    cb.checked = !cb.checked;
                }
            });
            featuresList.appendChild(div);
        });
    }

    // Handle Generation
    generateBtn.addEventListener('click', async () => {
        const enabledFeatures = {};
        featuresList.querySelectorAll('input[type="checkbox"]').forEach(cb => {
            enabledFeatures[cb.dataset.featureId] = cb.checked;
        });

        const payload = {
            name: document.getElementById('name').value,
            identifier: document.getElementById('identifier').value,
            publisher: document.getElementById('publisher').value,
            description: document.getElementById('description').value,
            template: templateSelect.value,
            enabled_features: enabledFeatures
        };

        if (!payload.name || !payload.identifier) {
            alert('Please fill in at least the Project Name and Identifier.');
            return;
        }

        generateBtn.disabled = true;
        generateBtn.textContent = 'Scaffolding...';

        try {
            const response = await fetch('/api/generate', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(payload)
            });

            const data = await response.json();
            
            if (response.ok) {
                statusPanel.classList.remove('hidden');
                statusText.textContent = `🧱 Success! ${data.message}`;
            } else {
                throw new Error(data.message || 'Generation failed');
            }
        } catch (err) {
            alert('Error: ' + err.message);
        } finally {
            generateBtn.disabled = false;
            generateBtn.textContent = 'Generate Project';
        }
    });
});
