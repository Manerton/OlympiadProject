<?php

namespace App\Services;

use App\Builder\ApplicationBuilder;
use App\Builder\EventBuilder;
use App\Builder\UserBuilder;
use App\Components\Dictionaries\ApplicationStatusDictionary;
use App\Repositories\ApplicationRepository;
use App\Repositories\EventRepository;
use App\Repositories\UserRepository;

class ApplicationService
{
    private ApplicationRepository $applicationRepository;
    private EventRepository $eventRepository;
    private UserRepository $userRepository;
    private ApplicationBuilder $applicationBuilder;
    private EventBuilder $eventBuilder;
    private UserBuilder $userBuilder;
    public function __construct(
        ApplicationRepository $applicationRepository,
        EventRepository $eventRepository,
        UserRepository $userRepository,
        ApplicationBuilder $applicationBuilder,
        EventBuilder $eventBuilder,
        UserBuilder $userBuilder
    )
    {
        $this->applicationRepository = $applicationRepository;
        $this->eventRepository = $eventRepository;
        $this->userRepository = $userRepository;
        $this->applicationBuilder = $applicationBuilder;
        $this->eventBuilder = $eventBuilder;
        $this->userBuilder = $userBuilder;
    }

    public function find($id)
    {
        $application = $this->applicationBuilder->build($this->applicationRepository->getByApiId($id));
        $event = $this->eventBuilder->build($this->eventRepository->getByApiId($application->event_id));
        $user = $this->userBuilder->build($this->userRepository->getByApiId($application->user_id));
        $this->applicationBuilder->buildUser($application, $user);
        $this->applicationBuilder->buildEvent($application, $event);
        return $application;
    }
    public function findAll($page)
    {
        $applications = [];
        $data = $this->applicationRepository->getByApiAll($page);
        foreach ($data as $item) {
            $application = $this->applicationBuilder->build($item);
            $event = $this->eventBuilder->build($this->eventRepository->getByApiId($application->event_id));
            $user = $this->userBuilder->build($this->userRepository->getByApiId($application->user_id));
            $this->applicationBuilder->buildUser($application, $user);
            $this->applicationBuilder->buildEvent($application, $event);
            $applications[] = $application;
        }
        return $applications;
    }
    public function findByEventId($eventId)
    {
        $applications = [];
        $data = $this->applicationRepository->getByEventId($eventId);
        foreach ($data as $item) {
            $application = $this->applicationBuilder->build($item);
            $event = $this->eventBuilder->build($this->eventRepository->getByApiId($eventId));
            $user = $this->userBuilder->build($this->userRepository->getByApiId($application->user_id));
            $this->applicationBuilder->buildUser($application, $user);
            $this->applicationBuilder->buildEvent($application, $event);
            $applications[] = $application;
        }
        return $applications;
    }
    public function confirmedApplications($applications)
    {
        $array = array_filter($applications, function ($application) {
            return $application->status == ApplicationStatusDictionary::APPROVED;
        });
        return $array;
    }
}
