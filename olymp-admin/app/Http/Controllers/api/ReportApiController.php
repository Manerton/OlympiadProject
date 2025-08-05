<?php

namespace App\Http\Controllers\api;

use App\Components\Dictionaries\SubjectDictionary;
use App\Components\ExcelCreator;
use App\Http\Controllers\Controller;
use App\Repositories\EventRepository;
use App\Services\ApplicationService;
use App\Services\EventService;
use App\Services\ReportService;

class ReportApiController extends Controller
{
    public EventRepository $eventRepository;
    public EventService $eventService;
    public ApplicationService $applicationService;
    public ReportService $reportService;
    public function __construct(
        EventRepository $eventRepository,
        EventService $eventService,
        ApplicationService $applicationService,
        ReportService $reportService
    )
    {
        $this->eventRepository = $eventRepository;
        $this->eventService = $eventService;
        $this->applicationService = $applicationService;
        $this->reportService = $reportService;
    }

    public function index(){
        $subjects = SubjectDictionary::getList();
        return response()->json([
            'subjects' => $subjects
        ]);
    }
    public function download($id){
        $data = $this->reportService->prepareData($id);
        return ExcelCreator::createList($id, $data);
    }
}
