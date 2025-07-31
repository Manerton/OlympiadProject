{{-- resources/views/participants/show.blade.php --}}
@extends('layouts.main')

@section('title', 'Просмотр участника: ' . $participant->userAPI->getFullFio())

@section('content')
    <div class="participant-view container mt-4">
        <h1>Просмотр участника: {{ $participant->userAPI->getFullFio() }}</h1>

        <p>
            <a href="{{ route('participant.edit', $participant->id) }}" class="btn btn-primary">Редактировать</a>

        <form action="{{ route('participant.delete', $participant->id) }}" method="POST" style="display:inline;">
            @csrf
            @method('DELETE')
            <button type="submit" class="btn btn-danger" onclick="return confirm('Вы уверены, что хотите удалить этого участника?')">
                Удалить
            </button>
        </form>
        </p>

        <table class="table table-striped table-bordered detail-view">
            <tbody>
            <tr>
                <th>ФИО</th>
                <td>{{ $participant->userAPI->getFullFio() }}</td>
            </tr>
            <tr>
                <th>Гражданство</th>
                <td>{{ $countries[$participant->citizenship] ?? '—' }}</td>
            </tr>
            <tr>
                <th>ОВЗ</th>
                <td>{{ $disabilities[$participant->disability] ?? '—' }}</td>
            </tr>
            <tr>
                <th>Обр. учреждение</th>
                <td>{{ $participant->schoolAPI->name ?? '—' }}</td>
            </tr>
            <tr>
                <th>Класс обучения</th>
                <td>{{ $classes[$participant->class] ?? ($participant->class . ' класс') }}</td>
            </tr>
            </tbody>
        </table>
    </div>
@endsection
