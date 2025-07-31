@extends('layouts.main')

@section('title', 'Просмотр образовательного учреждения ' . $model->name)

@section('breadcrumbs')
    @parent
    <li class="breadcrumb-item"><a href="{{ route('school.index') }}">Список обр. учреждений</a></li>
    <li class="breadcrumb-item active">Просмотр образовательного учреждения {{ $model->name }}</li>
@endsection

@section('content')
    <div class="school-view">
        <h1>Просмотр образовательного учреждения {{ $model->name }}</h1>

        <p>
            <a href="{{ route('school.edit', $model->id) }}" class="btn btn-primary">Редактировать</a>
        <form action="{{ route('school.delete', $model->id) }}" method="POST" style="display: inline-block;">
            @csrf
            @method('DELETE')
            <button type="submit" class="btn btn-danger"
                    onclick="return confirm('Вы уверены, что хотите удалить эту олимпиаду?')">
                Удалить
            </button>
        </form>
        </p>

        <table class="table table-striped table-bordered">
            <tbody>
            <tr>
                <th>Название образовательного учреждения</th>
                <td>{{ $model->name }}</td>
            </tr>
            <tr>
                <th>Регион</th>
                <td>{{ $regions[$model->region] }}</td>
            </tr>
            </tbody>
        </table>
    </div>
@endsection
