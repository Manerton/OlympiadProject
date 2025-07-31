<?php

namespace App\Http\Controllers;

use App\Components\Dictionaries\AttendanceDictionary;
use App\Components\Dictionaries\SubjectDictionary;
use App\Http\Requests\EventScoreRequest;
use App\Http\Requests\TaskRequest;
use App\Repositories\AttendanceRepository;
use App\Repositories\EventRepository;
use App\Repositories\TaskAttendanceRepository;
use App\Repositories\TaskRepository;
use App\Services\ApplicationService;
use App\Services\AttendanceService;
use App\Services\EventJuryService;
use App\Services\EventScoreService;
use App\Services\EventService;
use App\Services\TaskService;
use Illuminate\Http\Request;

class EventController extends Controller
{
    private EventService $eventService;
    private EventRepository $eventRepository;
    private ApplicationService $applicationService;
    private AttendanceRepository $attendanceRepository;
    private TaskRepository $taskRepository;
    private TaskAttendanceRepository $taskAttendanceRepository;
    private AttendanceService $attendanceService;
    private TaskService $taskService;
    private EventScoreService $eventScoreService;
    private EventJuryService $eventJuryService;
    public function __construct(
        EventService $eventService,
        EventRepository $eventRepository,
        ApplicationService $applicationService,
        AttendanceRepository $attendanceRepository,
        TaskRepository $taskRepository,
        TaskAttendanceRepository $taskAttendanceRepository,
        AttendanceService $attendanceService,
        TaskService $taskService,
        EventScoreService $eventScoreService,
        EventJuryService $eventJuryService
    )
    {
        $this->eventService = $eventService;
        $this->eventRepository = $eventRepository;
        $this->applicationService = $applicationService;
        $this->attendanceRepository = $attendanceRepository;
        $this->taskRepository = $taskRepository;
        $this->taskAttendanceRepository = $taskAttendanceRepository;
        $this->attendanceService = $attendanceService;
        $this->taskService = $taskService;
        $this->eventScoreService = $eventScoreService;
        $this->eventJuryService = $eventJuryService;
    }

    public function index($page = 1){
        $events = $this->eventService->findAll($page);
        $eventsAmount = $this->eventRepository->getCount();
        $subjects = SubjectDictionary::getList();
        return view('event/index', compact('events', 'eventsAmount', 'subjects'));
    }
    public function show($id){
        $event = $this->eventService->find($id);
        $eventJuries = $this->eventJuryService->findByEventId($event->id);
        $subjects = SubjectDictionary::getList();
        return view('event/show', compact('event', 'subjects', 'eventJuries'));
    }
    public function task($id)
    {
        $event = $this->eventService->find($id);
        $tasks = $this->taskRepository->getByEventId($id);
        return view('event/task', compact('event', 'tasks'));
    }
    public function attendance($id)
    {
        $event = $this->eventService->find($id);
        $applications = $this->applicationService->findByEventId($event->id);
        $applications = $this->applicationService->confirmedApplications($applications);
        $attendanceStatuses = AttendanceDictionary::getList();
        $attendances = $this->attendanceService->applicationFilter($applications);
        $data = $this->attendanceService->createTable($attendances);
        return view('event/attendance', compact('applications', 'event' , 'attendanceStatuses', 'data'));
    }
    public function point($id)
    {
        $event = $this->eventService->find($id);
        $applications = $this->applicationService->findByEventId($event->id);
        $applications = $this->applicationService->confirmedApplications($applications);
        $attendances = $this->attendanceService->attendanceFilter($this->attendanceService->applicationFilter($applications));
        $table = $this->attendanceService->createExtraTable($attendances);
        $tasks = $this->taskRepository->getByEventId($id);
        return view('event/point', compact('event', 'tasks', 'table'));
    }
    public function delete($id){
        return redirect()->route('event.index');
    }
    public function synchronize($id)
    {
        $eventApplications = $this->applicationService->findByEventId($id);
        $applications = $this->applicationService->confirmedApplications($eventApplications);
        $this->eventService->synchronize($applications, $eventApplications);
        return redirect()->route('event.attendance', ['id' => $id]);
    }
    public function addTask(TaskRequest $request, $id)
    {
        $data = $request->validated();
        $this->taskService->createTasks($data, $id);
        return redirect()->route('event.task', ['id' => $id]);
    }
    public function deleteTask($id)
    {
        $task = $this->taskRepository->get($id);
        $this->taskService->delete($task);
        return redirect()->route('event.task', ['id' => $task->event_id]);
    }
    public function changeAttendance(Request $request){
        try {
            $attendance = $this->attendanceRepository->get($request->attendance_id);
            if ($attendance) {
                $tasks = $this->taskRepository->getByEventId($request->eventId);
                if($request->status == AttendanceDictionary::NO_ATTENDANCE){
                    $taskAttendances = $this->taskAttendanceRepository->getByAttendanceId($attendance->id);
                    foreach ($taskAttendances as $taskAttendance) {
                        $this->taskAttendanceRepository->delete($taskAttendance);
                    }
                }
                else {
                    foreach ($tasks as $task) {
                        $this->taskAttendanceRepository->create($attendance->id, $task->id, 0);
                    }
                 }
                $attendance->status = $request->status;
                $this->attendanceRepository->save($attendance);
                return response()->json([
                    'success' => true,
                    'message' => 'Статус успешно обновлен'
                ]);
            }
            return response()->json([
                'success' => false,
                'message' => 'Запись посещения не найдена'
            ], 404);

        } catch (\Exception $e) {
            return response()->json([
                'success' => false,
                'message' => 'Ошибка сервера: ' . $e->getMessage()
            ], 500);
        }
    }
    public function changeScore(Request $request)
    {
        try {
            $taskAttendance = $this->taskAttendanceRepository->get($request->task_attendance_id);
            if ($taskAttendance) {
                $taskAttendance->points = $request->points;
                $this->taskAttendanceRepository->save($taskAttendance);
                return response()->json([
                    'success' => true,
                    'message' => 'Статус успешно обновлен'
                ]);
            }
            return response()->json([
                'success' => false,
                'message' => 'Запись посещения не найдена'
            ], 404);

        } catch (\Exception $e) {
            return response()->json([
                'success' => false,
                'message' => 'Ошибка сервера: ' . $e->getMessage()
            ], 500);
        }
    }
    public function prizeScore($id)
    {
        $event = $this->eventService->find($id);
        $eventScore = $this->eventScoreService->create($id);
        return view('event/prize-score', compact('event', 'eventScore'));
    }
    public function setPrizeScore(EventScoreRequest $request, $id)
    {
        $data = $request->validated();
        $this->eventScoreService->fill($id, $data);
        return redirect()->route('event.show', ['id' => $id]);
    }
}
