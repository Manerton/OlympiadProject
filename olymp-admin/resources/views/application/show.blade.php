@extends('layouts.main')

@section('title', 'Просмотр заявки')

@section('breadcrumbs')
    @parent
    <li class="breadcrumb-item"><a href="{{ route('application.index') }}">Список заявок</a></li>
    <li class="breadcrumb-item active">Просмотр заявки</li>
@endsection

@section('content')
    <div class="application-show">
        <h1>Детали заявки #{{ $application->code }}</h1>
        <div class="card">
            <div class="card-body">
                <div class="row">
                    <div class="col-md-6">
                        <h5>Основная информация</h5>
                        <p><strong>Код заявки:</strong> {{ $application->code }}</p>
                        <p><strong>Статус:</strong> {{ $statuses[$application->status] }}</p>
                    </div>
                    <div class="col-md-6">
                        <h5>Участник</h5>
                        <p><strong>ФИО:</strong> {{ $application->userAPI->getFullFio() }}</p>
                    </div>
                </div>

                <hr>

                <div class="row mt-3">
                    <div class="col-md-6">
                        <h5>Мероприятие</h5>
                        <p><strong>Название:</strong> {{ $application->eventAPI->name }}</p>
                        <p><strong>Предмет:</strong> {{ $application->eventAPI->subject }}</p>
                    </div>
                </div>
            </div>
        </div>

        <div class="mt-3">
            <a href="{{ route('application.index') }}" class="btn btn-secondary">Назад к списку</a>
            <a href="{{ route('application.edit', $application->id) }}" class="btn btn-warning">Редактировать</a>
        </div>
    </div>
@endsection
