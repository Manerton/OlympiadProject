<!DOCTYPE html>
<html lang="{{ str_replace('_', '-', app()->getLocale()) }}" class="h-100">
<head>
    <meta charset="{{ config('app.charset', 'utf-8') }}">
    <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">
    <meta name="description" content="{{ $meta_description ?? '' }}">
    <meta name="keywords" content="{{ $meta_keywords ?? '' }}">
    <meta name="csrf-token" content="{{ csrf_token() }}">
    <link rel="icon" href="{{ asset('favicon.ico') }}" type="image/x-icon">
    <title>{{ $title ?? config('app.name') }}</title>
    <link rel="stylesheet" href="/css/bootstrap.min.css">
    <script src="https://code.jquery.com/jquery-3.6.4.min.js"></script>
    @stack('head')
</head>

<body class="d-flex flex-column h-100">
<header id="header">
    <nav class="navbar navbar-expand-md navbar-dark bg-dark fixed-top">
        <div class="container">
            <a class="navbar-brand" href="{{ url('/') }}">
                {{ config('app.name') }}
            </a>
            <button class="navbar-toggler" type="button" data-bs-toggle="collapse" data-bs-target="#navbarMain">
                <span class="navbar-toggler-icon"></span>
            </button>
            <div class="collapse navbar-collapse" id="navbarMain">
                <ul class="navbar-nav me-auto">
                    <li class="nav-item"><a class="nav-link" href="{{ route('homepage') }}">Главная</a></li>
                    <li class="nav-item"><a class="nav-link" href="{{ route('application.index') }}">Заявки</a></li>
                    <li class="nav-item"><a class="nav-link" href="{{ route('school.index') }}">Обр. учреждения</a></li>
                    <li class="nav-item"><a class="nav-link" href="{{ route('user.index') }}">Пользователи</a></li>
                    <li class="nav-item"><a class="nav-link" href="{{ route('participant.index') }}">Участники</a></li>
                    <li class="nav-item"><a class="nav-link" href="{{ route('event.index') }}">Олимпиады</a></li>
                    <li class="nav-item"><a class="nav-link" href="{{ route('report.index') }}">Отчёты</a></li>
                </ul>

                @if (!Cookie::get('username'))
                    <div class="d-flex">
                        <a class="btn btn-link text-decoration-none" href="{{ route('login') }}">Login</a>
                    </div>
                @else
                    <form method="POST" action="{{ route('logout') }}" class="d-flex">
                        @csrf
                        <button type="submit" class="btn btn-link text-decoration-none">
                            Logout ({{ json_decode(Cookie::get('username'))->email }})
                        </button>
                    </form>
                @endif
            </div>
        </div>
    </nav>
</header>

<main id="main" class="flex-shrink-0" role="main" style="padding-top: 70px">
    <div class="container">
        @if (!empty($breadcrumbs))
        <nav aria-label="breadcrumb">
            <ol class="breadcrumb">
                @foreach ($breadcrumbs as $breadcrumb)
                <li class="breadcrumb-item">
                    @if (!empty($breadcrumb['url']))
                    <a href="{{ $breadcrumb['url'] }}">{{ $breadcrumb['label'] }}</a>
                    @else
                    {{ $breadcrumb['label'] }}
                    @endif
                </li>
                @endforeach
            </ol>
        </nav>
        @endif

        @if (session('alert'))
        <div class="alert alert-{{ session('alert.type', 'info') }}">
            {{ session('alert.message') }}
        </div>
        @endif
        @yield('content')
    </div>
</main>

<footer id="footer" class="mt-auto py-3 bg-light">
    <div class="container">
        <div class="row text-muted">
            <div class="col-md-6 text-center text-md-start">&copy; My Company {{ date('Y') }}</div>
            <div class="col-md-6 text-center text-md-end">
                <a href="https://laravel.com" class="text-muted">Powered by Laravel</a>
            </div>
        </div>
    </div>
</footer>

@stack('scripts')
</body>
</html>
