<?php

namespace App\Http\Controllers\api;

use App\Components\Dictionaries\NotificationTypeDictionary;
use App\Components\Dictionaries\StatusDictionary;
use App\Components\Dictionaries\ReasonParticipantDictionary;
use App\Components\Dictionaries\SubjectDictionary;
use App\Components\RabbitMQHelper;
use App\Http\Controllers\Controller;
use App\Http\Requests\ApplicationRequest;
use App\Repositories\ApplicationRepository;
use App\Services\ApplicationService;
use App\Services\EventService;
use App\Services\RabbitMQService;
use App\Services\UserService;

class ApplicationApiController extends Controller
{
    private RabbitMQService $rabbitMQService;
    private ApplicationService $applicationService;
    private ApplicationRepository $applicationRepository;
    private UserService $userService;
    private EventService $eventService;
    public function __construct(
        RabbitMQService $rabbitMQService,
        ApplicationService $applicationService,
        ApplicationRepository $applicationRepository,
        UserService $userService,
        EventService $eventService
    )
    {
        $this->rabbitMQService = $rabbitMQService;
        $this->applicationService = $applicationService;
        $this->applicationRepository = $applicationRepository;
        $this->userService = $userService;
        $this->eventService = $eventService;
    }

    public function index($page = 1){
        $applications = $this->applicationService->findAll($page);
        $applicationsAmount = $this->applicationRepository->getCount();
        $statuses = StatusDictionary::getList();
        $subjects = SubjectDictionary::getList();
        return response()->json([
            'applications' => collect($applications)->map(function($application) {
                return [
                    'id' => $application->id,
                    'user_id' => $application->user_id,
                    'reason' => $application->reason,
                    'status' => $application->status,
                    'code' => $application->code,
                    'userAPI' => [
                        'firstname' => $application->userAPI->firstname ?? null,
                        'surname' => $application->userAPI->surname ?? null,
                        'patronymic' => $application->userAPI->patronymic ?? null,
                    ],
                    'eventAPI' => [
                        'name' => $application->eventAPI->name ?? null,
                    ]
                ];
            })->toArray(),
            'applicationsAmount' => $applicationsAmount,
            'statuses' => $statuses,
            'subjects' => $subjects,
            'currentPage' => $page,
            'perPage' => 10
        ]);
    }
    public function create(){
        $statuses = StatusDictionary::getList();
        $users = $this->userService->findAll();
        $events = $this->eventService->findAll();
        return response()->json([
           'statuses' => $statuses,
           'users' => array($users),
           'events' => $events
        ]);
    }
    public function store(ApplicationRequest $request){
        $data = $request->validated();
        $this->rabbitMQService->publish(
            [RabbitMQHelper::APPLICATION_QUEUE_NAME],
            RabbitMQHelper::ADMIN_QUEUE_NAME,
            RabbitMQHelper::CREATE,
            RabbitMQHelper::APPLICATION_TABLE,
            array_diff_key($data, ['id' => null]),
        );
        return response()->json([]);
    }
    public function show($id){
        $application = $this->applicationService->find($id);
        $statuses = StatusDictionary::getList();
        return response()->json([
            'application' => [
                'id' => $application->id,
                'user_id' => $application->user_id,
                'reason' => $application->reason,
                'status' => $application->status,
                'code' => $application->code,
                'userAPI' => [
                    'firstname' => $application->userAPI->firstname ?? null,
                    'surname' => $application->userAPI->surname ?? null,
                    'patronymic' => $application->userAPI->patronymic ?? null,
                ],
                'eventAPI' => [
                    'name' => $application->eventAPI->name ?? null,
                ]
            ],
            'statuses' => $statuses,
        ]);
    }
    public function edit($id){
        $application = $this->applicationService->find($id);
        $statuses = StatusDictionary::getList();
        $users = $this->userService->findAll();
        $subjects = SubjectDictionary::getList();
        $reasons = ReasonParticipantDictionary::getList();
        $events = $this->eventService->findAll();
        return response()->json([
            'application' => [
                'id' => $application->id,
                'user_id' => $application->user_id,
                'reason' => $application->reason,
                'status' => $application->status,
                'code' => $application->code,
                'userAPI' => [
                    'firstname' => $application->userAPI->firstname ?? null,
                    'surname' => $application->userAPI->surname ?? null,
                    'patronymic' => $application->userAPI->patronymic ?? null,
                ],
                'eventAPI' => [
                    'name' => $application->eventAPI->name ?? null,
                ]
            ],
            'statuses' => $statuses,
            'users' => (array)$users,
            'subjects' => $subjects,
            'reasons' => $reasons,
            'events' => (array)$events,
        ]);
    }
    public function update(ApplicationRequest $request, $id){
        $data = $request->validated();
        $this->rabbitMQService->publish(
            [RabbitMQHelper::APPLICATION_QUEUE_NAME],
            RabbitMQHelper::ADMIN_QUEUE_NAME,
            RabbitMQHelper::UPDATE,
            RabbitMQHelper::APPLICATION_TABLE,
            array_diff_key($data, ['id' => null]),
            ['id' => $id]
        );
        return response()->json([]);
    }
    public function destroy($id){
        $this->rabbitMQService->publish(
            [RabbitMQHelper::APPLICATION_QUEUE_NAME],
            RabbitMQHelper::ADMIN_QUEUE_NAME,
            RabbitMQHelper::DELETE,
            RabbitMQHelper::APPLICATION_TABLE,
            [],
            ['id' => $id]
        );
        return response()->json([]);
    }
    public function confirm($id){
        $application = $this->applicationService->find($id);
        $this->rabbitMQService->publish(
            [RabbitMQHelper::APPLICATION_QUEUE_NAME],
            RabbitMQHelper::ADMIN_QUEUE_NAME,
            RabbitMQHelper::UPDATE,
            RabbitMQHelper::APPLICATION_TABLE,
            ['status' => StatusDictionary::APPROVED],
            ['id' => $id]
        );
        $this->rabbitMQService->publish(
            [RabbitMQHelper::NOTIFICATION_QUEUE_NAME],
            RabbitMQHelper::ADMIN_QUEUE_NAME,
            RabbitMQHelper::CREATE,
            RabbitMQHelper::NOTIFICATION_TABLE,
            [
                'user_id' => $application->user_id,
                'message' => NotificationTypeDictionary::APPLICATION_CONFIRM,
                'status' => NotificationTypeDictionary::ONLINE_NOTIFICATION
            ]
        );
        return response()->json([]);
    }
    public function reject($id){
        $application = $this->applicationService->find($id);
        $this->rabbitMQService->publish(
            [RabbitMQHelper::APPLICATION_QUEUE_NAME],
            RabbitMQHelper::ADMIN_QUEUE_NAME,
            RabbitMQHelper::UPDATE,
            RabbitMQHelper::APPLICATION_TABLE,
            ['status' => StatusDictionary::REJECTED],
            ['id' => $id]
        );
        $this->rabbitMQService->publish(
            [RabbitMQHelper::NOTIFICATION_QUEUE_NAME],
            RabbitMQHelper::ADMIN_QUEUE_NAME,
            RabbitMQHelper::CREATE,
            RabbitMQHelper::NOTIFICATION_TABLE,
            [
                'user_id' => $application->user_id,
                'message' => NotificationTypeDictionary::APPLICATION_REJECT,
                'status' => NotificationTypeDictionary::ONLINE_NOTIFICATION
            ]
        );
        return response()->json([]);
    }
}
