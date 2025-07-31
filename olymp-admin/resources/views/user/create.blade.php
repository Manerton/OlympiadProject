@extends('layouts.main')

@section('title', 'Создание нового пользователя')

@section('breadcrumbs')
    <li class="breadcrumb-item"><a href="{{ route('user.index') }}">Список пользователей</a></li>
    <li class="breadcrumb-item active">Создание нового пользователя</li>
@endsection

@section('content')
    <div class="user-form">
        <h1>Создание нового пользователя</h1>

        <form id="dynamic-form" method="POST" action="{{ route('user.store') }}">
            @csrf

            <div class="form-group">
                <label for="email">Email (логин)</label>
                <input type="email" class="form-control" id="email" name="email" maxlength="255" required>
                @error('email')<div class="text-danger">{{ $message }}</div>@enderror
            </div>

            <div class="form-group">
                <label for="password">Пароль</label>
                <input type="password" class="form-control" id="password" name="password" required>
                @error('password')<div class="text-danger">{{ $message }}</div>@enderror
            </div>

            <div class="form-group">
                <label for="surname">Фамилия</label>
                <input type="text" class="form-control" id="surname" name="surname" maxlength="255" required>
                @error('surname')<div class="text-danger">{{ $message }}</div>@enderror
            </div>

            <div class="form-group">
                <label for="firstname">Имя</label>
                <input type="text" class="form-control" id="firstname" name="firstname" maxlength="255" required>
                @error('firstname')<div class="text-danger">{{ $message }}</div>@enderror
            </div>

            <div class="form-group">
                <label for="patronymic">Отчество</label>
                <input type="text" class="form-control" id="patronymic" name="patronymic" maxlength="255">
                @error('patronymic')<div class="text-danger">{{ $message }}</div>@enderror
            </div>

            <div class="form-group">
                <label for="phone_number">Телефон</label>
                <input type="tel" class="form-control" id="phone_number" name="phone_number"
                       placeholder="+7 999 999-99-99">
                @error('phone_number')<div class="text-danger">{{ $message }}</div>@enderror
            </div>

            <div class="form-group">
                <label for="gender">Пол</label>
                <select class="form-control" id="gender" name="gender" required>
                    @foreach($genders as $key => $value)
                        <option value="{{ $key }}">{{ $value }}</option>
                    @endforeach
                </select>
                @error('gender')<div class="text-danger">{{ $message }}</div>@enderror
            </div>

            <div class="form-group">
                <label for="role">Роль</label>
                <select class="form-control" id="role" name="role" required>
                    @foreach($roles as $key => $value)
                        <option value="{{ $key }}">{{ $value }}</option>
                    @endforeach
                </select>
                @error('role')<div class="text-danger">{{ $message }}</div>@enderror
            </div>

            <div class="form-group">
                <label for="birthdate">Дата рождения</label>
                <input type="date" class="form-control" id="birthdate" name="birthdate" value="{{ old('birthdate') }}"
                       max="{{ date('Y-m-d') }}" required>
                @error('birthdate')<div class="text-danger">{{ $message }}</div>@enderror
            </div>

            <button type="submit" class="btn btn-primary">Сохранить</button>
        </form>
    </div>
@endsection
