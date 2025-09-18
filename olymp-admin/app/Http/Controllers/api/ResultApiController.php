<?php

namespace App\Http\Controllers\api;

use App\Components\ApiHelper;
use App\Components\Dictionaries\EventStatusDictionary;
use App\Http\Controllers\Controller;
use App\Repositories\AttendanceRepository;
use App\Services\ApplicationService;
use App\Services\EventService;

class ResultApiController extends Controller
{
    private AttendanceRepository $attendanceRepository;
    private ApplicationService $applicationService;
    private EventService $eventService;
    public function __construct(
        AttendanceRepository $attendanceRepository,
        ApplicationService $applicationService,
        EventService $eventService
    )
    {
        $this->attendanceRepository = $attendanceRepository;
        $this->applicationService = $applicationService;
        $this->eventService = $eventService;
    }
    public function resultByAttendance($id){
        $attendance = $this->attendanceRepository->get($id);
        return response()->json([
            ApiHelper::prepareResponse($attendance->taskAttendances)
        ]);
    }
    public function resultByUser($id){
        $applications = $this->applicationService->apiFindByUserId($id);
        $data = [];
        foreach ($applications as $application){
            $data[] = [
                'application' => $application,
                'attendance' =>  $this->attendanceRepository->getByApplicationId($application->id)[0],
                'result' =>  $this->attendanceRepository->getByApplicationId($application->id)[0]->taskAttendances
            ];
        }
        return response()->json(ApiHelper::prepareResponse($data));
    }
    public function resultByEventUser($eventId, $userId){
        $applications = $this->applicationService->apiFindByUserId($userId);
        $applications = array_filter($applications, function($application) use ($eventId){
            return $application->event_id == $eventId;
        });
        $data = [];
        foreach ($applications as $application){
            $attendance = $this->attendanceRepository->getByApplicationId($application->id)[0];
            foreach ($attendance->taskAttendances as $taskAttendance){
                $data[] = [
                    'task_id' => $taskAttendance->task_id,
                    'task_number' => $taskAttendance->task->number,
                    'points' => $taskAttendance->points,
                    'type' => $taskAttendance->task->type,
                ];
            }
        }
        return response()->json(ApiHelper::prepareResponse($data));
    }
    public function eventsByUser($id){
        $applications = $this->applicationService->apiFindByUserId($id);
        $events = [];
        foreach ($applications as $application){
            $attendance = $this->attendanceRepository->getByApplicationId($application->id)[0];
            if (isset($attendance->taskAttendances)){
                $events[] = $application->event_id;
            }
        }
        return response()->json(ApiHelper::prepareResponse($events));
        /*
            $applications = $this->applicationService->apiFindByUserId($id);
            $events = array_filter($this->eventService->findAll(), function($event){
                $event->conclude = EventStatusDictionary::CONCLUDE_OFF;
            });
            return response()->json(ApiHelper::prepareResponse(
                array_intersect(array_column($events, 'id'), array_column($applications, 'event_id'))
            ));
         */
    }
}
