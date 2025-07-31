<?php

namespace App\Services;

use App\Builder\ParticipantBuilder;
use App\Builder\SchoolBuilder;
use App\Builder\UserBuilder;
use App\Components\Dictionaries\RoleDictionary;
use App\Repositories\ParticipantRepository;
use App\Repositories\SchoolRepository;
use App\Repositories\UserRepository;
use Couchbase\Role;

class UserService
{
    private UserRepository $userRepository;
    private UserBuilder $userBuilder;
    private ParticipantRepository $participantRepository;
    private ParticipantBuilder $participantBuilder;
    private SchoolRepository $schoolRepository;
    private SchoolBuilder $schoolBuilder;
    public function __construct(
        UserRepository $userRepository,
        UserBuilder $userBuilder,
        ParticipantRepository $participantRepository,
        ParticipantBuilder $participantBuilder,
        SchoolRepository $schoolRepository,
        SchoolBuilder $schoolBuilder
    )
    {
        $this->userRepository = $userRepository;
        $this->userBuilder = $userBuilder;
        $this->participantRepository = $participantRepository;
        $this->participantBuilder = $participantBuilder;
        $this->schoolRepository = $schoolRepository;
        $this->schoolBuilder = $schoolBuilder;
    }

    public function find($id)
    {
        $user = $this->userBuilder->build($this->userRepository->getByApiId($id));
        if($user->role == RoleDictionary::PARTICIPANT) {
            $participant = $this->participantBuilder->build($this->participantRepository->getByApiUserId($id));
            $this->userBuilder->buildParticipant($user, $participant);
            $school = $this->schoolBuilder->build($this->schoolRepository->getByApiId($participant->school_id));
            $this->participantBuilder->buildSchool($participant, $school);
        }
        return $user;
    }
    public function findAll($page = NULL){
        $users = [];
        $data = $this->userRepository->getByApiAll($page);
        foreach ($data as $item) {
            $user = $this->userBuilder->build($item);
            if($user->role == RoleDictionary::PARTICIPANT) {
                $participant = $this->participantBuilder->build($this->participantRepository->getByApiUserId($user->id));
                $this->userBuilder->buildParticipant($user, $participant);
                $school = $this->schoolBuilder->build($this->schoolRepository->getByApiId($participant->school_id));
                $this->participantBuilder->buildSchool($participant, $school);
            }
            $users[] = $user;
        }
        return $users;
    }

}
