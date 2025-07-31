<?php

namespace App\Services;

use App\Builder\EventBuilder;
use App\Components\Dictionaries\ApplicationStatusDictionary;
use App\Repositories\AttendanceRepository;
use App\Repositories\EventRepository;

class EventService
{
    public AttendanceRepository $attendanceRepository;
    private EventRepository $eventRepository;
    private EventBuilder $eventBuilder;
    public function __construct(
        AttendanceRepository $attendanceRepository,
        EventRepository $eventRepository,
        EventBuilder $eventBuilder
    )
    {
        $this->attendanceRepository = $attendanceRepository;
        $this->eventRepository = $eventRepository;
        $this->eventBuilder = $eventBuilder;
    }

    public function find($id)
    {
        $event = $this->eventBuilder->build($this->eventRepository->getByApiId($id));
        return $event;
    }
    public function findAll($page = NULL)
    {
        $events = [];
        $data = $this->eventRepository->getByApiAll($page);
        foreach ($data as $event) {
            $events[] = $this->eventBuilder->build($event);
        }
        return $events;
    }
    public function synchronize($applications, $eventApplications){
        foreach ($applications as $application){
            if (count($this->attendanceRepository->getByApplicationId($application->id)) == 0 && $application->status == ApplicationStatusDictionary::APPROVED){
                $this->attendanceRepository->create($application->id);
            }
        }
        $currentAttendances = $this->attendanceRepository->getByApplicationsId(array_column($eventApplications, 'id'));
        $deletingIds = array_diff(array_column($currentAttendances->toArray(), 'application_id'), array_column($applications,'id'));
        foreach ($deletingIds as $id){
            $attendances = $this->attendanceRepository->getByApplicationId($id);
            foreach ($attendances as $attendance) {
                $this->attendanceRepository->delete($attendance);
            }
        }
    }
    public function filterBySubject($events, $subject)
    {
        return array_filter($events, function ($event) use ($subject) {
            return $event->subject == $subject;
        });
    }

}
