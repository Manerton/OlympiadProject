<div class="dynamic-form">
    <div id="dynamic-fields-wrapper">
        <div class="dynamic-field row align-items-end mb-2">
            @foreach($attributes as $attribute)
                @foreach($attribute as $item)
                    <div class="col-md-6">
                        <label class="form-label">{{ $item['label'] }}</label>
                        <input type="{{ $item['type'] }}"
                               class="form-control"
                               name="{{ $item['name'] }}[]" required>
                    </div>
                @endforeach
            @endforeach
            <div class="col-md-2">
                <button type="button" class="btn btn-danger remove-field">−</button>
            </div>
        </div>
    </div>
    <button type="button" id="add-field" class="btn btn-primary mb-3">+</button>
</div>

<script>
    document.addEventListener('DOMContentLoaded', function () {
        const addButton = document.getElementById('add-field');
        const wrapper = document.getElementById('dynamic-fields-wrapper');
        const fieldTemplate = wrapper.querySelector('.dynamic-field').cloneNode(true);

        addButton.addEventListener('click', function () {
            const newField = fieldTemplate.cloneNode(true);
            const inputs = newField.querySelectorAll('input');
            inputs.forEach(input => input.value = '');
            wrapper.appendChild(newField);
        });

        wrapper.addEventListener('click', function (e) {
            if (e.target.classList.contains('remove-field')) {
                const field = e.target.closest('.dynamic-field');
                if (wrapper.children.length > 1) {
                    field.remove();
                }
            }
        });
    });
</script>
