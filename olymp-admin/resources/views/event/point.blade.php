@extends('layouts.main')

@section('title', 'Выставление баллов')

@section('content')
    <a href="{{ route('event.show', $event->id) }}" class="btn btn-sm btn-primary">Перейти в карточку олимпиады</a>
    <div class="event-task">
        <h1>Выставление баллов</h1>
        <table class="table table-bordered table-striped mt-3">
            <thead>
            <tr>
                <th>ФИО участника</th>
                <th>Код участника</th>
                @foreach($tasks as $task)
                    <th>{{ $task->number }}</th>
                @endforeach
            </tr>
            </thead>
            <tbody>
            @foreach($table as $row)
                <tr>
                    <td>{{     $row['person']->getFullFio() }}</td>
                    <td>{{     $row['application']->code }}</td>
                    @foreach($row['taskAttendances'] as $counter => $taskAttendance)
                        <td> <input type="text"
                                    id="task_points-{{ $taskAttendance->id }}"
                                    name="task_points[{{ $taskAttendance->id }}]"
                                    value="{{ $taskAttendance->points ?? '' }}"
                                    maxlength="{{ $taskAttendance->task->max_points }}"
                                    oninput="this.value=this.value.replace(/[^0-9]/g,'');"></td>
                    @endforeach
                </tr>
            @endforeach
            </tbody>
        </table>
    </div>
    <script>
        $(document).ready(function() {
            // Обработчик изменения select
            $('input[id^="task_points-"]').change(function() {
                // Получаем ID из атрибута id (status-123 => 123)
                var token = $('meta[name="csrf-token"]').attr('content');
                var taskAttendanceId = this.id.split('-')[1];
                var points = $(this).val();
                console.log(taskAttendanceId, points);
                $.ajax({
                    url: '/event/change-score' , // Укажите ваш URL
                    type: 'POST',
                    headers: {
                        'X-CSRF-TOKEN': $('meta[name="csrf-token"]').attr('content')
                    },
                    data: {
                        _token: token,
                        task_attendance_id: taskAttendanceId,
                        points: points

                    },
                });
            });
        });
    </script>
@endsection
