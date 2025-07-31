@extends('layouts.main')

@section('title', 'Список пользователей')

@section('breadcrumbs')
    <li class="breadcrumb-item active">Список пользователей</li>
@endsection

@section('content')
    <div class="user-index">
        <p>Список пользователей</p>

        <a href="{{ route('user.create') }}" class="btn btn-success">Добавить пользователя</a>

        <table class="table">
            <thead>
            <tr>
                <th>#</th>
                <th>ФИО</th>
                <th>Эл.почта</th>
                <th>Номер телефона</th>
                <th>Действия</th>
            </tr>
            </thead>
            <tbody>
            @foreach($users as $index => $user)
                <tr>
                    <td>{{ $index + 1 }}</td>
                    <td>{{ $user->getFullFio() }}</td>
                    <td>{{ $user->email }}</td>
                    <td>{{ $user->phone_number }}</td>
                    <td>
                        <a href="{{ route('user.show', $user->id) }}" class="btn btn-sm btn-primary">Просмотр</a>
                        <a href="{{ route('user.edit', $user->id) }}" class="btn btn-sm btn-warning">Редактировать</a>
                        <form action="{{ route('user.delete', $user->id) }}" method="POST" style="display: inline-block;">
                            @csrf
                            @method('DELETE')
                            <button type="submit" class="btn btn-sm btn-danger" onclick="return confirm('Вы уверены, что хотите удалить этот элемент?')">Удалить</button>
                        </form>
                    </td>
                </tr>
            @endforeach
            </tbody>
        </table>
        <div class="pagination">
            @for($i = 0; $i <= $usersAmount / 10; $i++)
                <a href="{{ route('user.index', ['page' => $i + 1]) }}">{{ $i + 1 }}</a>
            @endfor
        </div>
    </div>
@endsection
