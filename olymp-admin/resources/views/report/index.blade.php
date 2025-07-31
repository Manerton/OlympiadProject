@extends('layouts.main')

@section('title', 'Список предметов')

@section('breadcrumbs')
    <li class="breadcrumb-item active">Список предметов</li>
@endsection

@section('content')
    <div class="report-index">
        <p>Список предметов</p>
        <table class="table">
            <thead>
            <tr>
                <th>#</th>
                <th>Предметы</th>
                <th>Действия</th>
            </tr>
            </thead>
            <tbody>
            @foreach($subjects as $index => $subject)
                <tr>
                    <td>{{ $index }}</td>
                    <td>{{ $subject }}</td>
                    <td>
                        <a href="{{ route('report.download', $index) }}" class="btn btn-sm btn-primary">Сформировать отчёт</a>
                    </td>
                </tr>
            @endforeach
            </tbody>
        </table>
    </div>
@endsection
