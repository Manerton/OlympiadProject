<?php

namespace App\Http\Controllers\api;

use App\Components\Dictionaries\ClassDictionary;
use App\Components\Dictionaries\CountryDictionary;
use App\Components\Dictionaries\DisabilityDictionary;
use App\Components\Dictionaries\GenderDictionary;
use App\Components\Dictionaries\RoleDictionary;
use App\Components\RabbitMQHelper;
use App\Http\Controllers\Controller;
use App\Http\Requests\ParticipantRequest;
use App\Repositories\ParticipantRepository;
use App\Services\ParticipantService;
use App\Services\RabbitMQService;
use App\Services\SchoolService;

class ParticipantApiController extends Controller
{
    private ParticipantRepository $participantRepository;
    private ParticipantService $participantService;
    private RabbitMQService $rabbitMQService;
    private SchoolService $schoolService;
    public function __construct(
        ParticipantService $participantService,
        ParticipantRepository $participantRepository,
        RabbitMQService $rabbitMQService,
        SchoolService $schoolService
    )
    {
        $this->participantService = $participantService;
        $this->participantRepository = $participantRepository;
        $this->rabbitMQService = $rabbitMQService;
        $this->schoolService = $schoolService;
    }

    public function index($page = 1)
    {
        $disabilities = DisabilityDictionary::getList();
        $countries = CountryDictionary::getList();
        $classes = ClassDictionary::getList();
        $participantsAmount = $this->participantRepository->getCount();
        $participants = $this->participantService->findAll($page);
        if (is_array($participants)) {
            $participants = collect($participants);
        }

        return response()->json([
            'disabilities' => $disabilities,
            'countries' => $countries,
            'classes' => $classes,
            'participants' => $participants->map(function($participant) {
                return [
                    'id' => $participant->id,
                    'user_id' => $participant->user_id,
                    'disability' => $participant->disability,
                    'citizenship' => $participant->citizenship,
                    'class' => $participant->class,
                    'school_id' => $participant->school_id,
                    'userAPI' => [
                        'firstname' => $participant->userAPI->firstname ?? null,
                        'surname' => $participant->userAPI->surname ?? null,
                        'patronymic' => $participant->userAPI->patronymic ?? null,
                    ],
                    'schoolAPI' => [
                        'name' => $participant->schoolAPI->name ?? null,
                    ]
                ];
            })->toArray(),
            'participantsAmount' => $participantsAmount,
            'currentPage' => $page,
            'perPage' => 10
        ]);
    }
    public function create(){
        $disabilities = DisabilityDictionary::getList();
        $countries = CountryDictionary::getList();
        $classes = ClassDictionary::getList();
        $schools = $this->schoolService->findAll();
        $roles = RoleDictionary::getList();
        $genders = GenderDictionary::getList();
        return response()->json([
            'disabilities' => $disabilities,
            'countries' => $countries,
            'classes' => $classes,
            'schools' => $schools,
            'roles' => $roles,
            'genders' => $genders
        ]);

    }
    public function store(ParticipantRequest $request){
        $data = $request->validated();
        $this->rabbitMQService->publish(
            [RabbitMQHelper::AUTH_QUEUE_NAME],
            RabbitMQHelper::ADMIN_QUEUE_NAME,
            RabbitMQHelper::CREATE,
            RabbitMQHelper::PARTICIPANT_TABLE,
            array_diff_key($data, ['id' => null]),
        );
        return response()->json([]);
    }
    public function show($id){
        $disabilities = DisabilityDictionary::getList();
        $countries = CountryDictionary::getList();
        $classes = ClassDictionary::getList();
        $participant = $this->participantService->find($id);
        return response()->json([
            'disabilities' => $disabilities,
            'countries' => $countries,
            'classes' => $classes,
            'participant' => [
                'id' => $participant->id,
                'user_id' => $participant->user_id,
                'disability' => $participant->disability,
                'citizenship' => $participant->citizenship,
                'class' => $participant->class,
                'school_id' => $participant->school_id,
                'userAPI' => [
                    'firstname' => $participant->userAPI->firstname ?? null,
                    'surname' => $participant->userAPI->surname ?? null,
                    'patronymic' => $participant->userAPI->patronymic ?? null,
                ],
                'schoolAPI' => [
                    'name' => $participant->schoolAPI->name ?? null,
                ]
            ]
        ]);
    }
    public function edit($id){
        $disabilities = DisabilityDictionary::getList();
        $countries = CountryDictionary::getList();
        $classes = ClassDictionary::getList();
        $roles = RoleDictionary::getList();
        $genders = GenderDictionary::getList();
        $schools = $this->schoolService->findAll();
        $participant = $this->participantService->find($id);
        return response()->json([
            'disabilities' => $disabilities,
            'countries' => $countries,
            'classes' => $classes,
            'participant' => [
                'id' => $participant->id,
                'user_id' => $participant->user_id,
                'disability' => $participant->disability,
                'citizenship' => $participant->citizenship,
                'class' => $participant->class,
                'school_id' => $participant->school_id,
                'userAPI' => [
                    'email' => $participant->userAPI->email ?? null,
                    'firstname' => $participant->userAPI->firstname ?? null,
                    'surname' => $participant->userAPI->surname ?? null,
                    'patronymic' => $participant->userAPI->patronymic ?? null,
                    'phone_number' => $participant->userAPI->phone_number ?? null,
                ],
                'schoolAPI' => [
                    'id' => $participant->schoolAPI->id ?? null,
                    'name' => $participant->schoolAPI->name ?? null,
                    'region' => $participant->schoolAPI->region ?? null,
                ]
            ],
            'schools' => array_map(fn($school) => (array)$school, $schools),
            'roles' => $roles,
            'genders' => $genders
        ]);
    }
    public function update(ParticipantRequest $request, $id){
        $data = $request->validated();
        $this->rabbitMQService->publish(
            [RabbitMQHelper::AUTH_QUEUE_NAME],
            RabbitMQHelper::ADMIN_QUEUE_NAME,
            RabbitMQHelper::UPDATE,
            RabbitMQHelper::PARTICIPANT_TABLE,
            array_diff_key($data, ['id' => null]),
            ['id' => $id]
        );

        return response()->json([]);
    }
    public function delete($id){
        $this->rabbitMQService->publish(
            [RabbitMQHelper::AUTH_QUEUE_NAME],
            RabbitMQHelper::ADMIN_QUEUE_NAME,
            RabbitMQHelper::DELETE,
            RabbitMQHelper::PARTICIPANT_TABLE,
            [],
            ['id' => $id]
        );
        return response()->json([]);
    }
}
