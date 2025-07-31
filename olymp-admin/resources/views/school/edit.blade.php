@extends('layouts.main')

@section('content')
    <form method="POST" action="{{ route('school.update', $model->id) }}" id="dynamic-form">
        @csrf
        <div class="form-group">
            <label for="name">Название образовательного учреждения</label>
            <input type="text" name="name" id="name"
                   class="form-control" maxlength="255"
                   value="{{ old('name', $model->name) }}">
            @error('name')
            <div class="text-danger">{{ $message }}</div>
            @enderror
        </div>

        <div class="form-group">
            <label for="region">Регион</label>
            <select name="region" id="region" class="form-control">
                @foreach($regions as $key => $region)
                    <option value="{{ $key }}" {{ old('region', $model->region) == $key ? 'selected' : '' }}>
                        {{ $region }}
                    </option>
                @endforeach
            </select>
            @error('region')
            <div class="text-danger">{{ $message }}</div>
            @enderror
        </div>

        <button type="submit" class="btn btn-primary">Обновить</button>
    </form>
@endsection
