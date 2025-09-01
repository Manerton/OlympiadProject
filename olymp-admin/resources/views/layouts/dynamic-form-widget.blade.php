<div class="dynamic-form">
    <div id="dynamic-fields-wrapper">
        <div class="dynamic-field row align-items-end mb-2">
            <div class="row g-4">
                @foreach($attributes as $attribute)
                    @foreach($attribute as $item)
                        <div class="col-12">
                            <div class="row align-items-center">
                                <div class="col-md-3">
                                    <label class="form-label fw-semibold text-dark">
                                        {{ $item['label'] }}
                                    </label>
                                </div>
                                <div class="col-md-6">
                                    @if($item['type'] != 'list')
                                        <input type="{{ $item['type'] }}"
                                               class="form-control"
                                               name="{{ $item['name'] }}[]"
                                               required>
                                    @else
                                        <select name="{{ $item['name'] }}[]"
                                                class="form-select">
                                            @foreach($item['elements'] as $index => $element)
                                                <option value="{{$index}}">{{$element}}</option>
                                            @endforeach
                                        </select>
                                    @endif
                                </div>
                                <div class="col-md-3">
                        <span class="text-muted small">
                            <i class="bi bi-asterisk text-danger"></i> Обязательно
                        </span>
                                </div>
                            </div>
                            <hr class="my-3">
                        </div>
                    @endforeach
                @endforeach
            </div>
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
