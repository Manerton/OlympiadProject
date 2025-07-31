<?php

namespace App\Http\Controllers;

use App\Components\Dictionaries\GenderDictionary;
use App\Components\Dictionaries\RoleDictionary;
use App\Components\RabbitMQHelper;
use App\Http\Requests\UserRequest;
use App\Repositories\UserRepository;
use App\Services\RabbitMQService;
use App\Services\UserService;

class UserController extends Controller
{
    private UserRepository $userRepository;
    private RabbitMQService $rabbitMQService;
    private UserService $userService;
    public function __construct(
        UserRepository $userRepository,
        RabbitMQService $rabbitMQService,
        UserService $userService
    )
    {
        $this->userRepository = $userRepository;
        $this->rabbitMQService = $rabbitMQService;
        $this->userService = $userService;
    }

    public function index($page = 1){
        $usersAmount = $this->userRepository->getCount();
        $users = $this->userService->findAll($page);
        return view('user/index', compact('users', 'usersAmount'));

    }
    public function create(){
        $roles = RoleDictionary::getList();
        $genders = GenderDictionary::getList();
        return view('user/create', compact('roles', 'genders'));
    }
    public function store(UserRequest $request){

        $data = $request->validated();
        $this->rabbitMQService->publish(
            [RabbitMQHelper::AUTH_QUEUE_NAME],
            RabbitMQHelper::QUEUE_NAME,
            RabbitMQHelper::CREATE,
            RabbitMQHelper::USER_TABLE,
            array_diff_key($data, ['id' => null]),
        );
        return redirect('/user/index');
    }
    public function show($id){
        $model = $this->userService->find($id);
        $roles = RoleDictionary::getList();
        $genders = GenderDictionary::getList();
        return view('user/show', compact('model', 'roles', 'genders'));
    }
    public function edit($id){
        $user = $this->userService->find($id);
        $roles = RoleDictionary::getList();
        $genders = GenderDictionary::getList();
        return view('user/edit', compact('user', 'roles', 'genders'));
    }
    public function update(UserRequest $request, $id){
        $data = $request->validated();
        $this->rabbitMQService->publish(
            [RabbitMQHelper::AUTH_QUEUE_NAME],
            RabbitMQHelper::QUEUE_NAME,
            RabbitMQHelper::UPDATE,
            RabbitMQHelper::USER_TABLE,
            array_diff_key($data, ['id' => null]),
            ['id' => $id]
        );
        return redirect('/user/index');
    }
    public function delete($id){
        $this->rabbitMQService->publish(
            [RabbitMQHelper::AUTH_QUEUE_NAME],
            RabbitMQHelper::QUEUE_NAME,
            RabbitMQHelper::DELETE,
            RabbitMQHelper::USER_TABLE,
            [],
            ['id' => $id]
        );
        return redirect('/user/index');
    }
}
