@extends('layouts.main')

@section('title', 'Список школ')

@section('breadcrumbs')
    @parent
    <li class="breadcrumb-item active"></li>
@endsection

@section('content')
    <div class="school-index">
        <p>Список школ</p>

        <a href="{{ route('school.create') }}" class="btn btn-success">Добавить обр. учреждение</a>

        <table class="table">
            <thead>
            <tr>
                <th>#</th>
                <th>Название образовательного учреждения</th>
                <th>Регион</th>
                <th>Действия</th>
            </tr>
            </thead>
            <tbody>
            @foreach($schools as $index => $school)
                <tr>
                    <td>{{ $index + 1 }}</td>
                    <td>{{ $school->name }}</td>
                    <td>{{ $regions[$school->region] }}</td>
                    <td>
                        <a href="{{ route('school.show', $school->id) }}" class="btn btn-sm btn-primary">Просмотр</a>
                        <a href="{{ route('school.edit', $school->id) }}" class="btn btn-sm btn-warning">Редактировать</a>
                        <form action="{{ route('school.delete', $school->id) }}" method="POST" style="display: inline-block;">
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
            @for($i = 0; $i <= $schoolsAmount / 10; $i++)
                <a href="{{ route('school.index', ['page' => $i + 1]) }}">{{ $i + 1 }}</a>
            @endfor
        </div>
    </div>
@endsection
