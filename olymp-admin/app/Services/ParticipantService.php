<?php

namespace App\Services;


use App\Builder\ParticipantBuilder;
use App\Builder\SchoolBuilder;
use App\Builder\UserBuilder;
use App\Repositories\ParticipantRepository;
use App\Repositories\SchoolRepository;
use App\Repositories\UserRepository;

class ParticipantService
{
    private ParticipantBuilder $participantBuilder;
    private ParticipantRepository $participantRepository;
    private UserBuilder $userBuilder;
    private UserRepository $userRepository;
    private SchoolBuilder $schoolBuilder;
    private SchoolRepository $schoolRepository;
    public function __construct(
        ParticipantBuilder $participantBuilder,
        ParticipantRepository $participantRepository,
        UserBuilder $userBuilder,
        UserRepository $userRepository,
        SchoolBuilder $schoolBuilder,
        SchoolRepository $schoolRepository
    )
    {
        $this->participantBuilder = $participantBuilder;
        $this->participantRepository = $participantRepository;
        $this->userBuilder = $userBuilder;
        $this->userRepository = $userRepository;
        $this->schoolBuilder = $schoolBuilder;
        $this->schoolRepository = $schoolRepository;
    }
    public function find($id)
    {
        $participant = $this->participantBuilder->build($this->participantRepository->getByApiId($id));
        $user = $this->userBuilder->build($this->userRepository->getByApiId($participant->user_id));
        $school = $this->schoolBuilder->build($this->schoolRepository->getByApiId($participant->school_id));
        $this->participantBuilder->buildUser($participant, $user);
        $this->participantBuilder->buildSchool($participant, $school);
        return $participant;
    }
    public function findAll($page){
        $participants = [];
        $data = $this->participantRepository->getByApiAll($page);
        foreach($data as $item){
            $participant = $this->participantBuilder->build($item);
            $user = $this->userBuilder->build($this->userRepository->getByApiId($participant->user_id));
            $school = $this->schoolBuilder->build($this->schoolRepository->getByApiId($participant->school_id));
            $this->participantBuilder->buildUser($participant, $user);
            $this->participantBuilder->buildSchool($participant, $school);
            $participants[] = $participant;
        }
        return $participants;
    }
}
