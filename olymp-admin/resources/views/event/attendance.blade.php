@extends('layouts.main')

@section('title', 'Явки')

@section('content')
    <a href="{{ route('event.show', $event->id) }}" class="btn btn-sm btn-primary">Перейти в карточку олимпиады</a>
    <table class="table table-bordered table-striped mt-3">
        <thead>
        <tr>
            <th>#</th>
            <th>ФИО участника</th>
            <th>Явка</th>
        </tr>
        </thead>
        <tbody>
        @foreach($data as $index => $item)
            <tr>
                <td>{{ $index + 1 }}</td>
                <td>{{ $item['person']->getFullFio() }}</td>
                <td>
                    <select name="status-{{$item['attendance']->id}}" id="status-{{$item['attendance']->id}}" class="form-control">
                        @foreach($attendanceStatuses as $key => $status)
                            <option value="{{ $key }}"  @if($item['attendance']->status == $key) selected @endif>{{ $status }}</option>
                        @endforeach
                    </select>
                </td>
            </tr>
        @endforeach
        </tbody>
    </table>
    <script>
        $(document).ready(function() {
            // Обработчик изменения select
            $('select[id^="status-"]').change(function() {
                // Получаем ID из атрибута id (status-123 => 123)
                var attendanceId = this.id.split('-')[1];
                var newStatus = $(this).val();
                var eventId = '{{$event->id}}';
                // CSRF токен для защиты Laravel
                var token = $('meta[name="csrf-token"]').attr('content');

                // Отправка AJAX запроса
                $.ajax({
                    url: '/event/change-attendance' , // Укажите ваш URL
                    type: 'POST',
                    headers: {
                        'X-CSRF-TOKEN': $('meta[name="csrf-token"]').attr('content')
                    },
                    data: {
                        _token: token,
                        attendance_id: attendanceId,
                        status: newStatus,
                        eventId: eventId
                    },
                });
            });
        });
    </script>
@endsection
