<?php

namespace App\Http\Controllers;

use App\Components\Dictionaries\RegionDictionary;
use App\Components\RabbitMQHelper;
use App\Http\Requests\SchoolRequest;
use App\Repositories\SchoolRepository;
use App\Services\RabbitMQService;
use App\Services\SchoolService;

class SchoolController extends Controller
{
    private SchoolRepository $schoolRepository;
    private SchoolService $schoolService;
    private RabbitMQService $rabbitMQService;
    public function __construct(
        SchoolRepository $schoolRepository,
        SchoolService $schoolService,
        RabbitMQService $rabbitMQService
    )
    {
        $this->schoolRepository = $schoolRepository;
        $this->schoolService = $schoolService;
        $this->rabbitMQService = $rabbitMQService;
    }

    public function index($page = 1){
        $schools = $this->schoolService->findAll($page);
        $schoolsAmount = $this->schoolRepository->getCount();
        $regions = RegionDictionary::REGIONS;
        return view('school/index')->with('schools', $schools)->with('schoolsAmount', $schoolsAmount)->with('regions', $regions);
    }
    public function create(){
        $regions = RegionDictionary::REGIONS;
        return view('school/create')->with('regions', $regions);
    }
    public function store(SchoolRequest $request){
        $data = $request->validated();
        $this->rabbitMQService->publish(
            [RabbitMQHelper::AUTH_QUEUE_NAME],
            RabbitMQHelper::QUEUE_NAME,
            RabbitMQHelper::CREATE,
            RabbitMQHelper::SCHOOL_TABLE,
            array_diff_key($data, ['id' => null]),
        );
        return redirect('/school/index');
    }
    public function show($id){
        $model = $this->schoolService->find($id);
        $regions = RegionDictionary::REGIONS;
        return view('school/show')->with('model', $model)->with('regions', $regions);
    }
    public function edit($id){
        $model = $this->schoolService->find($id);
        $regions = RegionDictionary::REGIONS;
        return view('school/edit')->with('model', $model)->with('regions', $regions);
    }
    public function update(SchoolRequest $request, $id){
        $data = $request->validated();
        $this->rabbitMQService->publish(
            [RabbitMQHelper::AUTH_QUEUE_NAME],
            RabbitMQHelper::QUEUE_NAME,
            RabbitMQHelper::UPDATE,
            RabbitMQHelper::SCHOOL_TABLE,
            array_diff_key($data, ['id' => null]),
            ['id' => $id]
        );
        return redirect('school/index');
    }
    public function delete($id){
        $this->rabbitMQService->publish(
            [RabbitMQHelper::AUTH_QUEUE_NAME],
            RabbitMQHelper::QUEUE_NAME,
            RabbitMQHelper::DELETE,
            RabbitMQHelper::SCHOOL_TABLE,
            [],
            ['id' => $id]
        );
        return redirect('/school/index');
    }
}
