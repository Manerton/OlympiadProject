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

    public function find($id, $token = null)
    {
        $application = $this->applicationBuilder->build($this->applicationRepository->getByApiId($id, $token));
        $event = $this->eventBuilder->build($this->eventRepository->getByApiId($application->event_id, $token));
        $user = $this->userBuilder->build($this->userRepository->getByApiId($application->user_id, $token));
        $this->applicationBuilder->buildUser($application, $user);
        $this->applicationBuilder->buildEvent($application, $event);
        return $application;
    }
    public function findAll($page, $token = null)
    {
        $applications = [];
        $data = $this->applicationRepository->getByApiAll($page ,10, $token);
        foreach ($data as $item) {
            $application = $this->applicationBuilder->build($item);
            $event = $this->eventBuilder->build($this->eventRepository->getByApiId($application->event_id, $token));
            $user = $this->userBuilder->build($this->userRepository->getByApiId($application->user_id, $token));
            $this->applicationBuilder->buildUser($application, $user);
            $this->applicationBuilder->buildEvent($application, $event);
            $applications[] = $application;
        }
        return $applications;
    }
    public function findByEventId($eventId, $token = null)
    {
        $applications = [];
        $data = $this->applicationRepository->getByEventId($eventId, $token);
        foreach ($data as $item) {
            $application = $this->applicationBuilder->build($item);
            $event = $this->eventBuilder->build($this->eventRepository->getByApiId($eventId, $token));
            $user = $this->userBuilder->build($this->userRepository->getByApiId($application->user_id, $token));
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
