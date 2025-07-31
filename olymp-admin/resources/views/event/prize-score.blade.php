@extends('layouts.main')

@section('title', $event->name)

@section('content')
    <div class="event-prize-score">
        <h1>Создание новой заявки</h1>
        <form action="{{ route('event.set-prize-score', ['id' => $event->id]) }}" method="POST">
            @csrf
            <div class="card mb-3">
                <div class="card-header">Минимальное количество баллов для получения статуса призёра</div>
                <div class="card-body">
                    <div class="form-group">
                        <label for="user_id">Баллы</label>
                        <input type = "text" name = "prize_score" value = {{$eventScore->prize_score}}>
                    </div>
                </div>
            </div>
            <div class="card mb-3">
                <div class="card-header">Минимальное количество баллов для получения статуса победителя</div>
                <div class="card-body">
                    <div class="form-group">
                        <label for="event_id">Баллы</label>
                        <input type = "text" name = "winner_score" value = {{$eventScore->winner_score}}>
                    </div>
                </div>
            </div>
            <div class="form-group">
                <button type="submit" class="btn btn-success">Определить количество баллов</button>
            </div>
        </form>
    </div>
@endsection
