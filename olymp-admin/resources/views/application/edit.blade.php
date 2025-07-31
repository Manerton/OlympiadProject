@extends('layouts.main')

@section('title', 'Редактирование заявки')

@section('breadcrumbs')
    @parent
    <li class="breadcrumb-item"><a href="{{ route('application.index') }}">Список заявок</a></li>
    <li class="breadcrumb-item active">Редактирование заявки</li>
@endsection

@section('content')
    <div class="application-edit">
        <h1>Редактирование заявки #{{ $application->code . ' ' . $application->userAPI->getFullFio()}}</h1>

        <form action="{{ route('application.update', $application->id) }}" method="POST">
            @csrf
            <div class="card mb-3">
                <div class="card-header">Основная информация</div>
                <div class="card-body">
                    <div class="form-group">
                        <label for="status">Статус заявки</label>
                        <select class="form-control" id="status" name="status" required>
                            <option value="">Выберите статус</option>
                            @foreach($statuses as $key => $status)
                                <option value="{{ $key }}" @if($application->status == $key) selected @endif>
                                    {{ $status }}
                                </option>
                            @endforeach
                        </select>
                    </div>
                </div>
            </div>

            <div class="card mb-3">
                <div class="card-header">Информация об участнике</div>
                <div class="card-body">
                    <div class="form-group">
                        <h2> {{ $application->userAPI->getFullFio()}}</h2>
                    </div>
                </div>
            </div>

            <div class="card mb-3">
                <div class="card-header">Информация о мероприятии</div>
                <div class="card-body">
                    <h2> {{ $application->eventAPI->name . ' ' . $subjects[$application->eventAPI->subject] }}</h2>
                </div>
            </div>
            <div class="card mb-3">
                <div class="card-header">Причина участия</div>
                <div class="card-body">
                    <div class="form-group">
                        <select class="form-control" id="reason" name="reason" required>
                            <option value="">Выберите статус</option>
                            @foreach($reasons as $key => $reason)
                                <option value="{{ $key }}" @if($application->reason == $key) selected @endif>
                                    {{ $reason }}
                                </option>
                            @endforeach
                        </select>
                    </div>
                </div>
            </div>
            <div class="form-group">
                <button type="submit" class="btn btn-primary">Сохранить изменения</button>
                <a href="{{ route('application.index') }}" class="btn btn-secondary">Отмена</a>
            </div>
        </form>
    </div>
@endsection
