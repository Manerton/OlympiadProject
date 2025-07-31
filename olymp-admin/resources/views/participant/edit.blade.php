{{-- resources/views/participants/edit.blade.php --}}
@extends('layouts.main')

@section('title', 'Редактировать участника')

@section('content')
    <div class="container mt-4">
        <h1>Редактировать участника</h1>

        <form action="{{ route('participant.update', $participant->id) }}" method="POST">
            @csrf
            <div class="form-group">
                <label for="email">Email (логин)</label>
                <input type="email" class="form-control" id="email" name="email" value="{{ old('email', $participant->userAPI->email) }}" maxlength="255" required>
                @error('email')<div class="text-danger">{{ $message }}</div>@enderror
            </div>

            <div class="form-group">
                <label for="password">Новый пароль (оставьте пустым, если не нужно менять)</label>
                <input type="password" class="form-control" id="password" name="password">
                @error('password')<div class="text-danger">{{ $message }}</div>@enderror
            </div>

            <div class="form-group">
                <label for="surname">Фамилия</label>
                <input type="text" class="form-control" id="surname" name="surname" value="{{ old('surname', $participant->userAPI->surname) }}" maxlength="255" required>
                @error('surname')<div class="text-danger">{{ $message }}</div>@enderror
            </div>

            <div class="form-group">
                <label for="firstname">Имя</label>
                <input type="text" class="form-control" id="firstname" name="firstname" value="{{ old('firstname', $participant->userAPI->firstname) }}" maxlength="255" required>
                @error('firstname')<div class="text-danger">{{ $message }}</div>@enderror
            </div>

            <div class="form-group">
                <label for="patronymic">Отчество</label>
                <input type="text" class="form-control" id="patronymic" name="patronymic" value="{{ old('patronymic', $participant->userAPI->patronymic) }}" maxlength="255">
                @error('patronymic')<div class="text-danger">{{ $message }}</div>@enderror
            </div>

            <div class="form-group">
                <label for="phone_number">Телефон</label>
                <input type="tel" class="form-control" id="phone_number" name="phone_number"
                       value="{{ old('phone_number', $participant->userAPI->phone_number) }}">
                @error('phone_number')<div class="text-danger">{{ $message }}</div>@enderror
            </div>

            <div class="form-group">
                <label for="gender">Пол</label>
                <select class="form-control" id="gender" name="gender" required>
                    @foreach($genders as $key => $value)
                        <option value="{{ $key }}" {{ old('gender', $participant->userAPI->gender) == $key ? 'selected' : '' }}>{{ $value }}</option>
                    @endforeach
                </select>
                @error('gender')<div class="text-danger">{{ $message }}</div>@enderror
            </div>

            <div class="form-group">
                <label for="role">Роль</label>
                <select class="form-control" id="role" name="role" required>
                    @foreach($roles as $key => $value)
                        <option value="{{ $key }}" {{ old('role', $participant->userAPI->role) == $key ? 'selected' : '' }}>{{ $value }}</option>
                    @endforeach
                </select>
                @error('role')<div class="text-danger">{{ $message }}</div>@enderror
            </div>

            <div class="form-group">
                <label for="birthdate">Дата рождения</label>
                <input type="date" class="form-control" id="birthdate" name="birthdate"
                       value="{{ old('birthdate', $participant->userAPI->birthdate) }}"
                       max="{{ date('Y-m-d') }}" required>
                @error('birthdate')<div class="text-danger">{{ $message }}</div>@enderror
            </div>
            <div class="mb-3">
                <label for="disability" class="form-label">Ограничения по здоровью</label>
                <select name="disability" id="disability" class="form-control">
                    @foreach($disabilities as $key => $label)
                        <option value="{{ $key }}" {{ $participant->disability == $key ? 'selected' : '' }}>{{ $label }}</option>
                    @endforeach
                </select>
            </div>

            <div class="mb-3">
                <label for="class" class="form-label">Класс обучения</label>
                <select name="class_number" id="class" class="form-control">
                    @foreach($classes as $key => $label)
                        <option value="{{ $key }}" {{ $participant->class == $key ? 'selected' : '' }}>{{ $label }}</option>
                    @endforeach
                </select>
            </div>

            <div class="mb-3">
                <label for="citizenship" class="form-label">Гражданство</label>
                <select name="citizenship" id="citizenship" class="form-control">
                    @foreach($countries as $key => $label)
                        <option value="{{ $key }}" {{ $participant->citizenship == $key ? 'selected' : '' }}>{{ $label }}</option>
                    @endforeach
                </select>
            </div>

            <div class="mb-3">
                <label for="school_id" class="form-label">Школа</label>
                <select name="school_id" id="school_id" class="form-control">
                    @foreach($schools as $school)
                        <option value="{{ $school->id }}" {{ $participant->school_id == $school->id ? 'selected' : '' }}>{{ $school->name }}</option>
                    @endforeach
                </select>
            </div>

            <button type="submit" class="btn btn-primary">Сохранить</button>
        </form>
    </div>
@endsection
