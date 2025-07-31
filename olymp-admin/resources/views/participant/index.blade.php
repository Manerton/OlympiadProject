@extends('layouts.main')

@section('title', 'Список участников деятельности')

@section('content')
    <div class="participant-index">
        <p>Список участников деятельности</p>

        <a href="{{ route('participant.create') }}" class="btn btn-success">Добавить участника деятельности</a>

        <table class="table table-bordered table-striped mt-3">
            <thead>
            <tr>
                <th>#</th>
                <th>ФИО</th>
                <th>Гражданство</th>
                <th>ОВЗ</th>
                <th>Обр. учреждение</th>
                <th>Класс обучения</th>
                <th>Действия</th>
            </tr>
            </thead>
            <tbody>
            @foreach($participants as $index => $participant)
                <tr>
                    <td>{{ $index + 1 }}</td>
                    <td>{{ $participant->userAPI->getFullFio() }}</td>
                    <td>{{ $countries[$participant->citizenship] ?? '—' }}</td>
                    <td>{{ $disabilities[$participant->disability] ?? '—' }}</td>
                    <td>{{ $participant->schoolAPI->name }}</td>
                    <td>{{ $participant->class }} класс</td>
                    <td>
                        <a href="{{ route('participant.show', $participant->id) }}" class="btn btn-sm btn-primary">Просмотр</a>
                        <a href="{{ route('participant.edit', $participant->id) }}" class="btn btn-sm btn-warning">Редактировать</a>
                        <form action="{{ route('participant.delete', $participant->id) }}" method="POST" style="display:inline;">
                            @csrf
                            @method('DELETE')
                            <button class="btn btn-sm btn-danger" onclick="return confirm('Вы уверены, что хотите удалить этот элемент?')">Удалить</button>
                        </form>
                    </td>
                </tr>
            @endforeach
            </tbody>
        </table>

        <div class="pagination">
            @for ($i = 1; $i <= ceil($participantsAmount / 10); $i++)
                <a href="{{ route('participant.index', ['page' => $i]) }}">{{ $i }}</a>
            @endfor
        </div>
    </div>
@endsection
