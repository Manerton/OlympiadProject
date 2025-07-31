<?php

namespace App\Services;

use App\Repositories\EventScoreRepository;
use App\Repositories\TaskAttendanceRepository;

class ReportService
{
    private ApplicationService $applicationService;
    private EventService $eventService;
    private AttendanceService $attendanceService;
    private UserService $userService;
    private TaskAttendanceRepository $taskAttendanceRepository;
    private EventScoreRepository $eventScoreRepository;
    public function __construct(
        ApplicationService $applicationService,
        EventService $eventService,
        AttendanceService $attendanceService,
        UserService $userService,
        TaskAttendanceRepository $taskAttendanceRepository,
        EventScoreRepository $eventScoreRepository
    )
    {
        $this->applicationService = $applicationService;
        $this->eventService = $eventService;
        $this->attendanceService = $attendanceService;
        $this->userService = $userService;
        $this->taskAttendanceRepository = $taskAttendanceRepository;
        $this->eventScoreRepository = $eventScoreRepository;
    }
    public function prepareData($id)
    {
        $data = [];
        $events = $this->eventService->filterBySubject($this->eventService->findAll(), $id);

        usort($events, function($a, $b) {
            return $a->class_number <=> $b->class_number; // Сортировка по возрастанию
        });
        foreach ($events as $event){
            $eventScore = $this->eventScoreRepository->getByEventId($event->id);
            $applications = $this->applicationService->findByEventId($event->id);
            $applications = $this->applicationService->confirmedApplications($applications);
            $attendances = $this->attendanceService->applicationFilter($applications);
            $participantData = [];
            foreach ($attendances as $attendance) {
                $taskAttendances = $this->taskAttendanceRepository->getByAttendancesId(array_column([$attendance], 'id'));
                $application = $this->applicationService->find($attendance->application_id);
                $person = $this->userService->find($application->user_id);
                $participantData[] = [
                    'attendance' => $attendance,
                    'application' => $application,
                    'person' => $person,
                    'taskAttendances' => $taskAttendances,
                ];
            }
            $data[] = [
                'event' => $event,
                'participants' => $participantData,
                'prizeScore' => $eventScore ? $eventScore->prize_score : -1,
                'winnerScore' => $eventScore ? $eventScore->winner_score : -1,
            ];
        }
        return $data;
    }
}
