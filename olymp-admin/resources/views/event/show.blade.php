@extends('layouts.main')

@section('title', $event->name)

@section('content')
    <div class="event-show">
        <h1>{{ $event->name }}</h1>

        <div class="card mt-4">
            <div class="card-body">
                <div class="row">
                    <div class="col-md-6">
                        <h5 class="card-title">Основная информация</h5>
                        <ul class="list-group list-group-flush">
                            <li class="list-group-item">
                                <strong>Предмет:</strong> {{ $subjects[$event->subject] }}
                            </li>
                            <li class="list-group-item">
                                <strong>Дата начала:</strong> {{ \Carbon\Carbon::parse($event->start_date)->format('d.m.Y H:i') }}
                            </li>
                            <li class="list-group-item">
                                <strong>Дата окончания:</strong> {{ \Carbon\Carbon::parse($event->end_date)->format('d.m.Y H:i') }}
                            </li>
                            <li class="list-group-item">
                                <strong>Возрастная категория:</strong> {{ $event->class_number }} класс
                            </li>
                            <li class="list-group-item">
                                <strong>Список жюри:</strong>
                                @foreach($eventJuries as $eventJury)
                                    <br>{{$eventJury->userAPI->getFullFio() }}
                                @endforeach
                            </li>
                            <li class="list-group-item">
                                <strong>Статус:</strong>
                                @if(now() < $event->start_date)
                                    <span class="badge bg-warning text-dark">Предстоящая</span>
                                @elseif(now() >= $event->start_date && now() <= $event->end_date)
                                    <span class="badge bg-success">Идет сейчас</span>
                                @else
                                    <span class="badge bg-secondary">Завершена</span>
                                @endif
                            </li>
                        </ul>
                    </div>
                </div>
            </div>

            <div class="card-footer">
                <div class="d-flex justify-content-between">
                    <a href="{{ route('event.index') }}" class="btn btn-secondary">
                        <i class="fas fa-arrow-left"></i> Назад к списку
                    </a>
                    <a href="{{ route('event.prize-score', $event->id) }}" class="btn btn-danger">
                        <i class="fas fa-arrow-left"></i> Перейти к определению баллов
                    </a>
                    <a href="{{ route('event.synchronize', $event->id) }}" class="btn btn-success">
                        <i class="fas fa-arrow-left"></i> Синхронизировать
                    </a>
                    <a href="{{ route('event.attendance', $event->id) }}" class="btn btn-primary">
                        <i class="fas fa-arrow-left"></i> Перейти к явкам
                    </a>
                    <a href="{{ route('event.task', $event->id) }}" class="btn btn-warning">
                        <i class="fas fa-arrow-left"></i> Перейти к заданиям
                    </a>
                    <a href="{{ route('event.point', $event->id) }}" class="btn btn-success">
                        <i class="fas fa-arrow-left"></i> Перейти к выставлению баллов
                    </a>
                    <div>
                        <form action="{{ route('event.delete', $event->id) }}" method="POST" class="d-inline">
                            @csrf
                            @method('DELETE')
                            <button type="submit" class="btn btn-danger" onclick="return confirm('Вы уверены, что хотите удалить эту олимпиаду?')">
                                <i class="fas fa-trash"></i> Удалить
                            </button>
                        </form>
                    </div>
                </div>
            </div>
        </div>
    </div>
@endsection
