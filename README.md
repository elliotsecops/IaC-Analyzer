**Análisis de Configuración de IaC**

## Resumen

El Análisis de Configuración de IaC es una herramienta diseñada para analizar archivos de configuración de Terraform HCL y realizar verificaciones de seguridad y optimización de costos. La herramienta identifica posibles problemas en las configuraciones de Terraform y los reporta con niveles de gravedad apropiados. Esto ayuda a mantener despliegues de infraestructura seguros y eficientes en términos de costos.
>>>>>>> bca134a (Updated IaC Analyzer with bug fix for severity counting and unique issue IDs)

## Contenido de la Tabla

1. [Características](#características)
2. [Instalación](#instalación)
3. [Uso](#uso)
4. [Verificaciones Implementadas](#verificaciones-implementadas)
5. [Estructura de Directorio](#estructura-de-directorio)
6. [Contribuciones](#contribuciones)
7. [Licencia](#licencia)

## Características

- **Verificaciones de Seguridad**: Identifica problemas de seguridad como el acceso SSH abierto, el acceso público a cubetas S3 y las instancias RDS sin cifrar.
- **Verificaciones de Optimización de Costos**: Identifica instancias EC2 de gran tamaño, volúmenes EBS grandes e IPs elásticos sin adjuntar.
- **Registros Verbosos**: Proporciona registros detallados para la depuración y análisis.
- **Interfaz de Línea de Comandos**: Fácil de usar con argumentos de línea de comandos.

## Instalación

### Requisitos Previos

- Python 3.7 o superior
- Pip (instalador de paquetes de Python)

### Pasos

1. **Clonar el Repositorio**:
   ```sh
   git clone https://github.com/elliotsecops/iac-analyzer.git
   cd iac-analyzer
   ```

2. **Crear un Entorno Virtual** (Opcional pero recomendado):
   ```sh
   python3 -m venv venv
   source venv/bin/activate  # En Windows, use `venv\Scripts\activate`
   ```

3. **Instalar Dependencias**:
   ```sh
   pip install -r requirements.txt
   ```

## Uso

### Interfaz de Línea de Comandos

La herramienta se puede ejecutar desde la línea de comandos con las siguientes opciones:

```sh
python main.py <path> [-v]
```

- `<path>`: Ruta al archivo de Terraform o directorio a analizar.
- `-v` o `--verbose`: Habilita la salida detallada para una salida más detallada.

### Ejemplos

1. **Analizar un Archivo de Terraform Individual**:
   ```sh
   python main.py tests/test_sample.tf
   ```

2. **Analizar un Directorio de Archivos de Terraform**:
   ```sh
   python main.py tests/
   ```

3. **Habilitar la Salida Verbosa**:
   ```sh
   python main.py tests/test_sample.tf -v
   ```

### Salida

La herramienta emite un resumen de los hallazgos, incluyendo el número de problemas encontrados y sus niveles de gravedad. Por ejemplo:

```
6 problemas encontrados (3 ALTA, 1 MEDIA, 2 BAJA, 1 INFORMACIÓN)
```

## Verificaciones Implementadas

### Verificaciones de Seguridad

1. **Acceso SSH Abierto**:
   - **Gravedad**: ALTA
   - **Descripción**: Identifica grupos de seguridad con acceso SSH abierto (puerto 22 abierto al 0.0.0.0/0).

2. **Acceso Público a Cubetas S3**:
   - **Gravedad**: MEDIA
   - **Descripción**: Identifica cubetas S3 con acceso de lectura público habilitado.

3. **Instancias RDS sin Cifrar**:
   - **Gravedad**: ALTA
   - **Descripción**: Identifica instancias RDS que no están cifradas.

### Verificaciones de Optimización de Costos

1. **Instancias EC2 de Gran Tamaño**:
   - **Gravedad**: INFORMACIÓN
   - **Descripción**: Identifica instancias EC2 que son de gran tamaño y sugiere reducir el tamaño.

2. **Volúmenes EBS Grandes**:
   - **Gravedad**: BAJA
   - **Descripción**: Identifica volúmenes EBS grandes y sugiere reducir el tamaño para reducir los costos.

3. **IPs Elásticos sin Adjuntar**:
   - **Gravedad**: BAJA
   - **Descripción**: Identifica IPs elásticos que no están adjuntos a ninguna instancia.

## Estructura de Directorio

```
iac-analyzer/
├── main.py
├── analyzer/
│   ├── __init__.py
│   ├── parser.py
│   ├── security.py
│   ├── cost.py
│   └── analyzer.py
├── tests/
│   ├── test_analyzer.py
│   ├── test_cost.tf
│   ├── test_sample.tf
│   └── test_security.tf
├── requirements.txt
└── README.md
```

- **`main.py`**: El script principal que maneja los argumentos de la línea de comandos, la salida y coordina el análisis.
- **`analyzer/`**: Contiene la lógica principal para analizar, verificar la seguridad, optimizar los costos y el módulo de análisis.
  - **`__init__.py`**: Archivo vacío para hacer que el directorio `analyzer` sea un paquete de Python.
  - **`parser.py`**: Analiza los archivos HCL de Terraform.
  - **`security.py`**: Implementa las verificaciones de seguridad.
  - **`cost.py`**: Implementa las verificaciones de optimización de costos.
  - **`analyzer.py`**: Coordina el análisis y la salida.
- **`tests/`**: Contiene los archivos de prueba y el script de prueba unitaria.
  - **`test_analyzer.py`**: Pruebas unitarias para el analizador.
  - **`test_cost.tf`**: Archivo de Terraform de ejemplo para la prueba de optimización de costos.
  - **`test_sample.tf`**: Archivo de Terraform de ejemplo para la prueba general.
  - **`test_security.tf`**: Archivo de Terraform de ejemplo para la prueba de seguridad.
- **`requirements.txt`**: Lista las dependencias de Python necesarias.
- **`README.md`**: Proporciona la descripción del proyecto, las instrucciones de instalación, los ejemplos de uso y una lista de las verificaciones implementadas.

## Contribuciones

Las contribuciones son bienvenidas. Por favor, siga estos pasos:

1. Haga un fork del repositorio.
2. Cree una nueva rama (`git checkout -b feature-branch`).
3. Realice sus cambios (`git commit -am 'Add some feature'`).
4. Empuje a la rama (`git push origin feature-branch`).
5. Cree una nueva pull request.

## Licencia

Este proyecto está licenciado bajo la Licencia MIT. Consulte el archivo [LICENSE](LICENSE) para obtener más detalles.

---

# IaC Configuration Analyzer

## Overview

The IaC Configuration Analyzer is a tool designed to parse Terraform HCL files and perform security and cost optimization checks. The tool identifies potential issues in Terraform configurations and reports them with appropriate severity levels. This helps in maintaining secure and cost-effective infrastructure deployments.

## Table of Contents

1. [Features](#features)
2. [Installation](#installation)
3. [Usage](#usage)
4. [Implemented Checks](#implemented-checks)
5. [Directory Structure](#directory-structure)
6. [Contributing](#contributing)
7. [License](#license)

## Features

- **Security Checks**: Identifies security issues such as open SSH access, public S3 bucket access, and unencrypted RDS instances.
- **Cost Optimization Checks**: Identifies oversized EC2 instances, large EBS volumes, and unattached Elastic IPs.
- **Verbose Logging**: Provides detailed logs for debugging and analysis.
- **Command-Line Interface**: Easy to use with command-line arguments.

## Installation

### Prerequisites

- Python 3.7 or higher
- Pip (Python package installer)

### Steps

1. **Clone the Repository**:
   ```sh
   git clone https://github.com/elliotsecops/iac-analyzer.git
   cd iac-analyzer
   ```

2. **Create a Virtual Environment** (Optional but recommended):
   ```sh
   python3 -m venv venv
   source venv/bin/activate  # On Windows use `venv\Scripts\activate`
   ```

3. **Install Dependencies**:
   ```sh
   pip install -r requirements.txt
   ```

## Usage

### Command-Line Interface

The tool can be run from the command line with the following options:

```sh
python main.py <path> [-v]
```

- `<path>`: Path to the Terraform file or directory to analyze.
- `-v` or `--verbose`: Enable verbose logging for detailed output.

### Examples

1. **Analyze a Single Terraform File**:
   ```sh
   python main.py tests/test_sample.tf
   ```

2. **Analyze a Directory of Terraform Files**:
   ```sh
   python main.py tests/
   ```

3. **Enable Verbose Logging**:
   ```sh
   python main.py tests/test_sample.tf -v
   ```

### Output

The tool outputs a summary of findings, including the number of issues found and their severity levels. For example:

```
6 issues found (3 HIGH, 1 MEDIUM, 2 LOW, 1 INFO)
```

## Implemented Checks

### Security Checks

1. **Open SSH Access**:
   - **Severity**: HIGH
   - **Description**: Identifies security groups with open SSH access (port 22 open to 0.0.0.0/0).

2. **Public S3 Bucket Access**:
   - **Severity**: MEDIUM
   - **Description**: Identifies S3 buckets with public read access enabled.

3. **Unencrypted RDS Instances**:
   - **Severity**: HIGH
   - **Description**: Identifies RDS instances that are not encrypted.

### Cost Optimization Checks

1. **Oversized EC2 Instances**:
   - **Severity**: INFO
   - **Description**: Identifies EC2 instances that are oversized and suggests downsizing.

2. **Large EBS Volumes**:
   - **Severity**: LOW
   - **Description**: Identifies large EBS volumes and suggests resizing to reduce cost.

3. **Unattached Elastic IPs**:
   - **Severity**: LOW
   - **Description**: Identifies Elastic IPs that are not attached to any instance.

## Directory Structure

```
iac-analyzer/
├── main.py
├── analyzer/
│   ├── __init__.py
│   ├── parser.py
│   ├── security.py
│   ├── cost.py
│   └── analyzer.py
├── tests/
│   ├── test_analyzer.py
│   ├── test_cost.tf
│   ├── test_sample.tf
│   └── test_security.tf
├── requirements.txt
└── README.md
```

- **`main.py`**: The main script that handles command-line arguments, logging, and orchestrates the analysis.
- **`analyzer/`**: Contains the core logic for parsing, security checks, cost optimization checks, and the analyzer module.
  - **`__init__.py`**: Empty file to make the `analyzer` directory a Python package.
  - **`parser.py`**: Parses Terraform HCL files.
  - **`security.py`**: Implements security checks.
  - **`cost.py`**: Implements cost optimization checks.
  - **`analyzer.py`**: Coordinates parsing and analysis.
- **`tests/`**: Contains test files and the unit test script.
  - **`test_analyzer.py`**: Unit tests for the analyzer.
  - **`test_cost.tf`**: Sample Terraform file for cost optimization testing.
  - **`test_sample.tf`**: Sample Terraform file for comprehensive testing.
  - **`test_security.tf`**: Sample Terraform file for security testing.
- **`requirements.txt`**: Lists the required Python dependencies.
- **`README.md`**: Provides project description, installation instructions, usage examples, and a list of implemented checks.

## Contributing

Contributions are welcome! Please follow these steps:

1. Fork the repository.
2. Create a new branch (`git checkout -b feature-branch`).
3. Make your changes and commit them (`git commit -am 'Add some feature'`).
4. Push to the branch (`git push origin feature-branch`).
5. Create a new Pull Request.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
