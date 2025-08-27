@extends('layouts.main')

@section('title', 'Сервис отправки почты')

@section('content')
    <div class="mail-index">
        <form method="POST" action="{{ route('mail.send') }}">
            @csrf
            <div class="mb-3">
                <label for="email" class="form-label">Email <span style="color:red">*</span></label>
                <input type="email" id="email" name="email" class="form-control" required>
                @error('email')
                <div class="text-danger">{{ $message }}</div>
                @enderror
            </div>
            <div class="mb-3">
                <label for="message" class="form-label">Сообщение <span style="color:red">*</span></label>
                <textarea id="message" name="message" rows="4" class="form-control" required></textarea>
                @error('message')
                <div class="text-danger">{{ $message }}</div>
                @enderror
            </div>
            <button type="submit" class="btn btn-primary">Отправить</button>
        </form>
    </div>
@endsection
