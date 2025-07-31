@extends('layouts.main')

@section('title', 'Создание новой заявки')

@section('breadcrumbs')
    @parent
    <li class="breadcrumb-item"><a href="{{ route('application.index') }}">Список заявок</a></li>
    <li class="breadcrumb-item active">Создание заявки</li>
@endsection

@section('content')
    <div class="application-create">
        <h1>Создание новой заявки</h1>
        <form action="{{ route('application.store') }}" method="POST">
            @csrf
            <div class="card mb-3">
                <div class="card-header">Выбор участника</div>
                <div class="card-body">
                    <div class="form-group">
                        <label for="user_id">Участник</label>
                        <select class="form-control" id="user_id" name="user_id" required>
                            <option value="">Выберите участника</option>
                            @foreach($users as $user)
                                <option value="{{ $user->id }}">{{ $user->getFullFio() }}</option>
                            @endforeach
                        </select>
                    </div>
                </div>
            </div>
            <div class="card mb-3">
                <div class="card-header">Выбор мероприятия</div>
                <div class="card-body">
                    <div class="form-group">
                        <label for="event_id">Мероприятие</label>
                        <select class="form-control" id="event_id" name="event_id" required>
                            <option value="">Выберите мероприятие</option>
                            @foreach($events as $event)
                                <option value="{{ $event->id }}">{{ $event->name }}</option>
                            @endforeach
                        </select>
                    </div>
                </div>
            </div>
            <div class="card mb-3">
                <div class="card-header">Статус</div>
                <div class="card-body">
                    <div class="form-group">
                        <option value="">Выберите статус</option>
                        @foreach($statuses as $index => $status)
                            <option value="{{ $index + 1 }}">{{ $status }}</option>
                        @endforeach
                    </div>
                </div>
            </div>
            <div class="form-group">
                <button type="submit" class="btn btn-success">Создать заявку</button>
                <a href="{{ route('application.index') }}" class="btn btn-secondary">Отмена</a>
            </div>
        </form>
    </div>
@endsection
