@extends('layouts.main')

@section('title', 'Список заявок')

@section('breadcrumbs')
    @parent
    <li class="breadcrumb-item active"></li>
@endsection

@section('content')
    <div class="application-index">
        <p>Список заявок</p>
        <table class="table">
            <thead>
            <tr>
                <th>#</th>
                <th>Код заявки</th>
                <th>ФИО участника</th>
                <th>Предмет</th>
                <th>Статус заявки</th>
                <th>Изменения статуса заявки</th>
                <th>Действия</th>
            </tr>
            </thead>
            <tbody>
            @foreach($applications as $index => $application)
                <tr>
                    <td>{{ $index + 1 }}</td>
                    <td>{{ $application->code }}</td>
                    <td>{{ $application->userAPI->getFullFio() }}</td>
                    <td>{{ $application->eventAPI->name . ' ' . $subjects[$application->eventAPI->subject] }}</td>
                    <td>{{ $statuses[$application->status] }}</td>
                    <td>
                        @if($application->status == \App\Components\Dictionaries\ApplicationStatusDictionary::AWAITING)
                            <form action="{{ route('application.confirm', $application->id) }}" method="POST" style="display: inline-block;">
                                @csrf
                                <button type="submit" class="btn btn-sm btn-success">Подтвердить заявку</button>
                            </form>
                            <form action="{{ route('application.reject', $application->id) }}" method="POST" style="display: inline-block;">
                                @csrf
                                <button type="submit" class="btn btn-sm btn-danger">Отклонить заявку</button>
                            </form>
                        @endif
                    </td>
                    <td>
                        <a href="{{ route('application.show', $application->id) }}" class="btn btn-sm btn-primary">Просмотр</a>
                        <a href="{{ route('application.edit', $application->id) }}" class="btn btn-sm btn-warning">Редактировать</a>
                        <form action="{{ route('application.delete', $application->id) }}" method="POST" style="display: inline-block;">
                            @csrf
                            @method('DELETE')
                            <button type="submit" class="btn btn-sm btn-danger" onclick="return confirm('Вы уверены, что хотите удалить этот элемент?')">Удалить</button>
                        </form>
                    </td>
                </tr>
            @endforeach
            </tbody>
        </table>

        <div class="pagination">
            @for($i = 0; $i <= $applicationsAmount / 10; $i++)
                <a href="{{ route('application.index', ['page' => $i + 1]) }}">{{ $i + 1 }}</a>
            @endfor
        </div>
    </div>
@endsection
