@extends('layouts.main')

@section('title', 'Просмотр пользователя ' . $model->getFullFio())

@section('breadcrumbs')
    <li class="breadcrumb-item"><a href="{{ route('user.index') }}">Список пользователей</a></li>
    <li class="breadcrumb-item active">Просмотр пользователя {{ $model->getFullFio() }}</li>
@endsection

@section('content')
    <div class="user-view">
        <h1>Просмотр пользователя {{ $model->getFullFio() }}</h1>

        <p>
            <a href="{{ route('user.edit', $model->id) }}" class="btn btn-primary">Редактировать</a>
        <form action="{{ route('user.delete', $model->id) }}" method="POST" style="display: inline-block;">
            @csrf
            @method('DELETE')
            <button type="submit" class="btn btn-danger" onclick="return confirm('Вы уверены, что хотите удалить этого пользователя?')">
                Удалить
            </button>
        </form>
        </p>

        <table class="table table-striped table-bordered">
            <tbody>
            <tr>
                <th>ФИО</th>
                <td>{{ $model->getFullFio() }}</td>
            </tr>
            <tr>
                <th>Email</th>
                <td>{{ $model->email }}</td>
            </tr>
            <tr>
                <th>Номер телефона</th>
                <td>{{ $model->phone_number }}</td>
            </tr>
            <tr>
                <th>Дата рождения</th>
                <td>{{ $model->birthdate }}</td>
            </tr>
            <tr>
                <th>Роль</th>
                <td>{{ $roles[$model->role] ?? 'Не указано' }}</td>
            </tr>
            <tr>
                <th>Пол</th>
                <td>{{ $genders[$model->gender] ?? 'Не указано' }}</td>
            </tr>
            <tr>
                <th>Дата рождения</th>
                <td>{{ $model->birthdate }}</td>
            </tr>
            </tbody>
        </table>
    </div>
@endsection
