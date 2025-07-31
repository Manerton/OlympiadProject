@php use App\Components\DynamicFormWidget; @endphp
@extends('layouts.main')

@section('title', 'Задания')

@section('content')
    <a href="{{ route('event.show', $event->id) }}" class="btn btn-sm btn-primary">Перейти в карточку олимпиады</a>

    <div class="event-task">
        <h1>Список заданий</h1>
        <form method="POST" action="{{ route('event.add-task', $event->id) }}">
            @csrf
            @php
                $widget = new DynamicFormWidget();
                $widget->attributes = [
                    'attributes' => [
                       [
                           'name' => 'number',
                           'label' => 'Номер задания',
                           'type' => 'text'
                       ],
                       [
                           'name' => 'point',
                           'label' => 'Количество баллов',
                           'type' => 'text'
                       ]
                    ]
                ];
                echo $widget->render();
            @endphp
            <button type="submit" class="btn btn-primary">Сохранить</button>
        </form>

        <table class="table table-bordered table-striped mt-3">
            <thead>
            <tr>
                <th>#</th>
                <th>Номер задания</th>
                <th>Максимальное количество баллов</th>
                <th>Действия</th>
            </tr>
            </thead>
            <tbody>
            @foreach($tasks as $index => $task)
                <tr>
                    <td>{{ $index + 1 }}</td>
                    <td>{{ $task->number }}</td>
                    <td>{{ $task->max_points }}</td>
                    <td>
                        <form action="{{ route('event.delete-task', $task->id) }}" method="POST" style="display:inline;">
                            @csrf
                            @method('DELETE')
                            <button class="btn btn-sm btn-danger" onclick="return confirm('Вы уверены, что хотите удалить этот элемент?')">Удалить</button>
                        </form>
                    </td>
                </tr>
            @endforeach
            </tbody>
        </table>
    </div>

@endsection
