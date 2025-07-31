@extends('layouts.main')

@section('title', 'Олимпиады')

@section('content')
    <div class="event-index">
        <p>Список олимпиад</p>

        <table class="table table-bordered table-striped mt-3">
            <thead>
            <tr>
                <th>#</th>
                <th>Название олимпиады</th>
                <th>Предмет</th>
                <th>Дата начала олимпиады</th>
                <th>Дата окончания олимпиады</th>
                <th>Возрастная категория</th>
                <th>Действия</th>
            </tr>
            </thead>
            <tbody>
            @foreach($events as $index => $event)
                <tr>
                    <td>{{ $index + 1 }}</td>
                    <td>{{ $event->name }}</td>
                    <td>{{ $subjects[$event->subject] }}</td>
                    <td>{{ $event->start_date }}</td>
                    <td>{{ $event->end_date }}</td>
                    <td>{{ $event->class_number }} класс</td>
                    <td>
                        <a href="{{ route('event.show', $event->id) }}" class="btn btn-sm btn-primary">Просмотр</a>
                        <form action="{{ route('event.delete', $event->id) }}" method="POST" style="display:inline;">
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
            @for ($i = 1; $i <= ceil($eventsAmount / 10); $i++)
                <a href="{{ route('event.index', ['page' => $i]) }}">{{ $i }}</a>
            @endfor
        </div>
    </div>
@endsection
